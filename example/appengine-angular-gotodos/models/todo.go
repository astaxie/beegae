// EXAMPLE FROM: https://github.com/GoogleCloudPlatform/appengine-angular-gotodos
//
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package models

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func DefaultTodoList(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "TodoList", "default", 0, nil)
}

type Todo struct {
	Id      int64     `json:"id" datastore:"-"`
	Text    string    `json:"text" datastore:",noindex"`
	Done    bool      `json:"done"`
	Created time.Time `json:"created"`
}

func (t *Todo) key(c context.Context) *datastore.Key {
	if t.Id == 0 {
		t.Created = time.Now()
		return datastore.NewIncompleteKey(c, "Todo", DefaultTodoList(c))
	}
	return datastore.NewKey(c, "Todo", "", t.Id, DefaultTodoList(c))
}

func (t *Todo) Save(c context.Context) (*Todo, error) {
	k, err := datastore.Put(c, t.key(c), t)
	if err != nil {
		return nil, err
	}
	t.Id = k.IntID()
	return t, nil
}
