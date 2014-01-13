// +build appengine

package session

import (
	"net/http"
	"sync"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

var appenginepvdr = &AppEngineProvider{}

type AppEngineSessionStore struct {
	c           appengine.Context
	sid         string
	lock        sync.RWMutex
	dirty       bool
	maxlifetime int64
	bss_entity  *BeegoSessionStore
	values      map[interface{}]interface{}
}

type BeegoSessionStore struct {
	SessionData  []byte `datastore:",noindex"`
	SessionStart time.Time
}

func (st *AppEngineSessionStore) Set(key, value interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values[key] = value
	st.dirty = true
	//st.updatestore()
	return nil
}

func (st *AppEngineSessionStore) Get(key interface{}) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	if v, ok := st.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *AppEngineSessionStore) Delete(key interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	delete(st.values, key)
	st.dirty = true
	//st.updatestore()
	return nil
}

func (st *AppEngineSessionStore) Flush() error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values = make(map[interface{}]interface{})
	st.dirty = true
	//st.updatestore()
	return nil
}

func (st *AppEngineSessionStore) SessionID() string {
	return st.sid
}

func (st *AppEngineSessionStore) updatestore() {
	b, err := encodeGob(st.values)
	if err != nil {
		st.c.Errorf("error encoding session data: %v", err)
		return
	}

	done := make(chan bool, 2)

	if st.bss_entity == nil {
		st.bss_entity = &BeegoSessionStore{SessionStart: time.Now()}
	}

	st.bss_entity.SessionData = b

	go func() {
		k := datastore.NewKey(st.c, "BeegoSessionStore", st.sid, 0, nil)
		if _, ds_err := datastore.Put(st.c, k, st.bss_entity); ds_err != nil {
			st.c.Errorf("error saving session data to datastore: %v", ds_err)
		}
		done <- true
	}()

	go func() {
		mem_err := memcache.Set(st.c, &memcache.Item{
			Key:        st.sid,
			Value:      st.bss_entity.SessionData,
			Expiration: (time.Duration(st.maxlifetime) * time.Second),
		})
		if mem_err != nil {
			st.c.Errorf("error saving session data to memcache: %v", mem_err)
		}
		done <- true
	}()

	_, _ = <-done, <-done
}

// SessionRelease will update the data of a session and reset its
// expiration time
func (st *AppEngineSessionStore) SessionRelease(w http.ResponseWriter) {
	//Always expected to be called to save session data
	st.bss_entity.SessionStart = time.Now()
	st.updatestore()
}

type AppEngineProvider struct {
	maxlifetime int64
}

func (mp *AppEngineProvider) SessionInit(gclifetime int64, config string) error {
	mp.maxlifetime = gclifetime
	return nil
}

func (mp *AppEngineProvider) getsession(sid string, c appengine.Context) *BeegoSessionStore {
	in_cache := false
	e := new(BeegoSessionStore)
	if item, err := memcache.Get(c, sid); err == nil {
		in_cache = true
		e.SessionData = item.Value
		e.SessionStart = time.Now() // Do we care about accuracy here?
	} else if err != memcache.ErrCacheMiss {
		c.Errorf("error getting session data from memcache: %v", err)
	}

	if !in_cache {
		k := datastore.NewKey(c, "BeegoSessionStore", sid, 0, nil)
		if ds_err := datastore.Get(c, k, e); ds_err != nil {
			e.SessionStart = time.Now()
			if ds_err != datastore.ErrNoSuchEntity {
				c.Errorf("error getting session data from datastore: %v", ds_err)
			}
		}
	}
	return e
}

func (mp *AppEngineProvider) SessionExist(sid string, c appengine.Context) bool {
	k := datastore.NewKey(c, "BeegoSessionStore", sid, 0, nil)
	e := new(BeegoSessionStore)
	if ds_err := datastore.Get(c, k, e); ds_err == datastore.ErrNoSuchEntity {
		return false
	} else if ds_err != nil {
		c.Errorf("error while checking existence of session data from datastore: %v", ds_err)
		return false
	} else {
		// Don't depend on GC to clean expired sessions
		return (time.Duration(mp.maxlifetime) * time.Second) > time.Since(e.SessionStart)
	}
}

func (mp *AppEngineProvider) SessionRead(sid string, c appengine.Context) (SessionStore, error) {
	e := mp.getsession(sid, c)
	var kv = make(map[interface{}]interface{})

	if len(e.SessionData) != 0 {
		decoded_gob, err := decodeGob(e.SessionData)
		if err != nil {
			return nil, err
		}
		kv = decoded_gob
	}
	rs := &AppEngineSessionStore{c: c, sid: sid, values: kv, maxlifetime: mp.maxlifetime, dirty: false, bss_entity: e}
	return rs, nil
}

func (mp *AppEngineProvider) SessionRegenerate(oldsid, sid string, c appengine.Context) (SessionStore, error) {
	e := mp.getsession(sid, c)
	var kv = make(map[interface{}]interface{})

	if len(e.SessionData) != 0 {
		decoded_gob, err := decodeGob(e.SessionData)
		if err != nil {
			return nil, err
		}
		kv = decoded_gob
	}
	rs := &AppEngineSessionStore{c: c, sid: sid, values: kv, maxlifetime: mp.maxlifetime, dirty: false, bss_entity: e}
	return rs, nil
}

func (mp *AppEngineProvider) SessionDestroy(sid string, c appengine.Context) error {
	done := make(chan bool, 2)

	go func() {
		k := datastore.NewKey(c, "BeegoSessionStore", sid, 0, nil)
		if ds_err := datastore.Delete(c, k); ds_err != nil {
			c.Errorf("error deleting session data from datastore: %v", ds_err)
		}
		done <- true
	}()

	go func() {
		mem_err := memcache.Delete(c, sid)
		if mem_err != nil {
			c.Errorf("error deleting session data from memcache: %v", mem_err)
		}
		done <- true
	}()

	_, _ = <-done, <-done
	return nil
}

func (mp *AppEngineProvider) SessionGC(c appengine.Context) {
	q := datastore.NewQuery("BeegoSessionStore").Filter("SessionStart <", time.Now().Unix()-mp.maxlifetime).KeysOnly()

	keys, err := q.GetAll(c, nil)
	if err != nil {
		c.Errorf("error querying session data from datastore: %v", err)
	}

	for _, key := range keys {
		mp.SessionDestroy(key.StringID(), c)
	}
	return
}

func (mp *AppEngineProvider) SessionAll() int {
	// Unable to Implement given Sessions API
	return 0
	/*
		total, err := datastore.NewQuery("BeegoSessionStore").KeysOnly().Count(c)
		if err != nil {
			return 0
		}
		return total
	*/
}

func init() {
	Register("appengine", appenginepvdr)
}
