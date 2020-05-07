package models

import (
	//"fmt"
	"time"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct{
	//Id int64 `orm:"pk";auto`
	Username string 		`orm:"pk"`
	Password string 
	LastLogout time.Time 	`orm:"type(datetime);null"`
}

func NewUser(user *User)(error){
	o := orm.NewOrm()
	_,err := o.Insert(user)
	return err
}

func GetUser(user *User)(error){
	o := orm.NewOrm()
	return o.Read(user)
}

func SetUser(user *User)(int64, error){
	o := orm.NewOrm()
	num, err := o.Update(user)
	return num, err
}

func RemoveUser(user *User)(int64, error){
	o := orm.NewOrm()
	num, err := o.Delete(user)
	return num, err
}

func AllUser()(*[]*User, error){
	var users []*User
	o := orm.NewOrm()
	if _,err := o.QueryTable("User").All(&users); err != nil{
		return nil, err
	}
	return &users, nil
}

func CheckPassword(username string, password string)bool{
	var users []*User
	o := orm.NewOrm()
	num,err := o.QueryTable("User").Filter("Username", username).Filter("Password", password).All(&users)
	if err != nil{
		return false
	}
	return (num == 1)
}