package utils

import(
	"github.com/astaxie/beego"
)

func Error(errs...interface{}){
	beego.Error(errs)
}

func Info(errs...interface{}){
	beego.Info(errs)
}