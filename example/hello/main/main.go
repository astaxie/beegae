package main

import (
	"github.com/astaxie/beegae"
	"github.com/astaxie/beegae/example/hello/controllers"
)

func init() {
	beegae.Router("/", &controllers.MainController{})
	beegae.Run()
}
