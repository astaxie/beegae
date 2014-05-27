beegae session
==============

This is based off of the original [session module](https://github.com/astaxie/beego/tree/master/session) as part of the beego project and can be used as part of the beegae project or as a standalone session manager. [read more here](http://beego.me/docs/mvc/controller/session.md)

This includes (as of now) only a single session store capable of working on AppEngine (`SessionProvider = "appengine"`)

A few gotchas:

1. There is no automatic garbage collection! You will have to create a cron job and a custom handler to periodically call on the garbage collection functions.
2. `SessionAll` will always return 0. `Count` queries are limited to 1000 entities and so we cannot reliably get a count. As such, this function was not implemented.
3. A few methods deviate from the original beego API specification. Specifically, an `appengine.Context` object is a new parameter for `SessionExist`, `SessionRead`, `SessionRegenerate`, `SessionDestroy`, and `SessionGC`. If you are using beegae or use the session manager provided, you do not have to worry about these details.
4. `GetProvider` was not implemented (this should have little to no impact)

Example Garbage Collection using **beegae**:

First, create a new controller:

```go
package controllers

import "github.com/astaxie/beegae"

type GCController struct {
	beegae.Controller
}

func (this *GCController) Get() {
	beegae.GlobalSessions.GC(this.AppEngineCtx)
}
```

Second, register your controller to a URL Path in your applications `init` function:

```go
func init() {
	// Register other routers/handlers here
	// ...

	// Register new handler for sessiong garbage collection
	beegae.Router("/_session_gc", &controllers.GCController)

	beegae.Run()
}
```

Finally, add an entry to your cron.yaml file:

```yaml
cron:
- description: daily session garbage collection
  url: /_session_gc
  schedule: every day 00:00
```

You can also add security to this (and any) URL by requiring an Admin login for the URL in your app.yaml:

```yaml
handlers:
- url: /_session_gc
  login: admin
  script: _go_app

- url: /.*
  script: _go_app
```
