package main

import (
	"controllers"
	"github.com/astaxie/beegae"
)

func init() {
	beegae.Router("/", &controllers.MainController{})
	beegae.Run()
}
