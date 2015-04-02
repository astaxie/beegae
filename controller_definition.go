// +build appengine !appenginevm

package beegae

import (
	"appengine"
	"github.com/astaxie/beegae/context"
	"github.com/astaxie/beegae/session"
)

// Controller defines some basic http request handler operations, such as
// http context, template and view, session and xsrf.
type Controller struct {
	Ctx            *context.Context
	Data           map[interface{}]interface{}
	controllerName string
	actionName     string
	TplNames       string
	Layout         string
	LayoutSections map[string]string // the key is the section name and the value is the template name
	TplExt         string
	_xsrf_token    string
	gotofunc       string
	CruSession     session.SessionStore
	XSRFExpire     int
	AppController  interface{}
	AppEngineCtx   appengine.Context
	EnableRender   bool
	EnableXSRF     bool
	methodMapping  map[string]func() //method:routertree
}

// Init generates default values of controller operations.
func (c *Controller) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Layout = ""
	c.TplNames = ""
	c.controllerName = controllerName
	c.actionName = actionName
	c.Ctx = ctx
	c.TplExt = "tpl"
	c.AppController = app
	c.EnableRender = true
	c.EnableXSRF = true
	c.Data = ctx.Input.Data
	c.AppEngineCtx = appengine.NewContext(ctx.Request)
	c.methodMapping = make(map[string]func())
}
