package main

import (
	"github.com/astaxie/beegae"
	"hello/controllers"
)

func init() {
	beegae.Router("/", &controllers.MainController{})
	beegae.Run()
}
