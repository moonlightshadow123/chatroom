package models

import (
	//"fmt"
	"time"
	"strconv"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Message struct{
	Id int64 		`orm:"pk;auto`
	Type int 
	Content string 	//`orm:"null"`
	Related string 	`orm:"null"`
	Date time.Time 	`orm:"auto_now_add;type(datetime)"`
	User *User 		`orm:"rel(fk)"`
}

func NewMsg(msg *Message)(error){
	o := orm.NewOrm()
	_,err := o.Insert(msg)
	return err
}

func GetMsg(msg *Message)(error){
	o := orm.NewOrm()
	return o.Read(msg)
}

func RemoveMsg(msg *Message)(int64, error){
	o := orm.NewOrm()
	num, err := o.Delete(msg)
	return num, err
}

func RecallMsg(msg *Message)(int64, error){
	var msgs []*Message
	o := orm.NewOrm()
	related, err := strconv.Atoi(msg.Related)
	if err!=nil{
		return 0, err
	}
	num1,err1 := o.QueryTable("Message").Filter("Id", related).Filter("User__Username", msg.User.Username).All(&msgs)
	if num1 != 1 || err1 != nil{
		return num1, err1
	}
	oldmsg := msgs[0]
	num2,err2 := RemoveMsg(oldmsg)
	return num2, err2
}

func GetHistMsg(t time.Time, num int)(*[]*Message, error){
    o := orm.NewOrm()
    var msgs []*Message
	_, err := o.QueryTable("Message").Filter("Date__lt", t).OrderBy("-Date").Limit(num).All(&msgs)
	if err != nil{
		return nil, err
	}
	return &msgs, nil 
	//fmt.Printf("Returned Rows Num: %s, %s", num, err)
}