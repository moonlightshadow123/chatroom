package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	_ "test_proj/chatRoom/routers"
	"test_proj/chatRoom/controllers"
)

//自动建表
func createTable() {
	name := "default"                          //数据库别名
	force := false                             //不强制建数据库
	verbose := true                            //打印建表过程
	err := orm.RunSyncdb(name, force, verbose) //建表
	if err != nil {
		beego.Error(err)
	}
}

func init() {
	// DB
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	createTable()
	// Session
	controllers.Init()
}

func main() {
	beego.Run()
}

