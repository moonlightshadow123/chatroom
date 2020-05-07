package controllers

import (
	"test_proj/chatRoom/ws"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get(){
	c.TplName = "index.html"
}


func Init(){
	man := ws.GetInst()
	man.Run()
}