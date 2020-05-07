package models
import (
	//"fmt"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Message))
}