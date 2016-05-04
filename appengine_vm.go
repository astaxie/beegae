// +build !appengine

package beegae

import (
	"net/http"

	"google.golang.org/appengine"
)

// Run beego application.
func (app *App) Run() {
	http.Handle("/", app.Handlers)
	appengine.Main()
}
