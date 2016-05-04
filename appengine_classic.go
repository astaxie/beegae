// +build appengine

package beegae

import "net/http"

// Run beego application.
func (app *App) Run() {
	http.Handle("/", app.Handlers)
}
