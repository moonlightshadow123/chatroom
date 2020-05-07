package routers

import (
	"github.com/astaxie/beego"
	"test_proj/chatRoom/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/chatRoom", &controllers.ServerController{})
	beego.Router("/logout", &controllers.ServerController{}, "get:Logout")
	beego.Router("/chatRoom/WS", &controllers.ServerController{}, "get:WS")
	beego.Router("/hist/:stamp:string/:num:int", &controllers.ServerController{}, "get:Hist")
}
