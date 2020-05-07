package controllers

import (
	"time"
	"strconv"
	//"net/http"
	//"encoding/json"
	"test_proj/chatRoom/ws"
	"test_proj/chatRoom/mymsg"
	"test_proj/chatRoom/utils"
	"test_proj/chatRoom/models"
	"github.com/astaxie/beego"
	//"github.com/gorilla/websocket"
	//"github.com/astaxie/beego/orm"
	//"github.com/astaxie/beego/session"
)

var(
	secret = "chatRoom"
)

type Message = mymsg.Message

type ServerController struct {
	beego.Controller
}

func (c *ServerController)Prepare(){
	if c.Ctx.Request.Method == "GET"{
		username := c.check()
		if username == ""{
		    beego.Error("Prepare Error.")
		    c.Ctx.Redirect(302, "/")
		}
    }
}

func (c *ServerController)check()string{
	username, _ := c.GetSecureCookie(secret, "username")
	return username
}

func getUserMap()(*map[string]bool, error){
	//var users *[]*User
	users, err:=models.AllUser()
	if err != nil{
		utils.Error("Get All User Error!", err)
		return nil, err
	}
	umap := make(map[string]bool)
	for _,user := range *users{
		umap[user.Username] = false
	}
	man := ws.GetInst()
	for name, _ := range man.Clients{
		umap[name] = true
	}
	return &umap, nil
}

func (c *ServerController)getChatHtml(username string){
	umap, err := getUserMap()
	if err != nil{
		utils.Error("Get All user Error!", err)
	}
	c.Data["umap"] = umap
	c.Data["stamp"] = time.Now().Unix()
	beego.Info("get name:"+username+", and send to chatRoom.html")
	c.Data["name"] = username
	c.TplName = "chatRoom.html"
}

func (c *ServerController)Post(){
	username := c.Ctx.Request.Form.Get("username")
	password := c.Ctx.Request.Form.Get("password")
	if !models.CheckPassword(username, password){
		beego.Error("Username or Passowrd Error!")
		c.Redirect("/", 302)
	}
	c.SetSecureCookie(secret, "username", username)
	c.getChatHtml(username)
}

func (c *ServerController)Get(){
	username := c.check()
	c.getChatHtml(username)
}

func (c *ServerController)Logout(){
	//c.GetSecureCookie(secret, "username") 
	c.SetSecureCookie(secret, "username", "")
	c.Redirect("/", 302)
}

type HistMsgs struct{
	MsgList []*Message `json:"msglist"`
	Stamp string `json:"stamp"`
}

func (c *ServerController)Hist(){
	num := c.Ctx.Input.Param(":num")
	stamp := c.Ctx.Input.Param(":stamp")
	number, ierr := strconv.Atoi(num)
	if ierr != nil{
		utils.Error("Int Number Paser Error. num =",num, "error =", ierr)
	}
	tstamp, errt := strconv.ParseInt(stamp, 10, 64)
    if errt != nil {
        utils.Error("Timestamp Paser Error. stamp =",stamp, "error =", ierr)
    }
    t := time.Unix(tstamp, 0)
	dbmsgs,err := models.GetHistMsg(t, number)
	if err != nil{
		utils.Error("Get History Msg Error!", err)
	}
	histmsgs := HistMsgs{MsgList:make([]*Message,0, number)}
	for _,dbmsg := range *dbmsgs{
		msg := mymsg.DbMsgToMsg(dbmsg)
		histmsgs.MsgList = append(histmsgs.MsgList, msg)
		histmsgs.Stamp = strconv.FormatInt(dbmsg.Date.Unix(), 10)
		//histmsgs.MsgList[i] = msg
	}
	c.Data["json"] = &histmsgs
	c.ServeJSON()
}

// 用于与用户间的websocket连接(chatRoom.html发送来的websocket请求)
func (c *ServerController) WS() {
	name := c.check()
	client, err := ws.NewClient(c.Ctx.ResponseWriter, c.Ctx.Request, name)
	if err != nil{
		utils.Error(err)
		return
	}
	man := ws.GetInst()
	// 如果用户列表中有该用户
	oldclient, has := man.Clients[name]
	if has {
		oldclient.Close()
	}
	man.Join <- client
	utils.Info("user:", client.Name, "websocket connect success!")

	// 当函数返回时，将该用户加入退出通道，并断开用户连接
	defer func() {
		man.Leave <- client
		client.Close()
		//c.relaseSess()
	}()

	// 由于WebSocket一旦连接，便可以保持长时间通讯，则该接口函数可以一直运行下去，直到连接断开
	for {
		// 读取消息。如果连接断开，则会返回错误
       	msg, err := client.OnFromWeb()
		if err == nil{
			man.Msgs <- msg
		}else{
			utils.Error(err)
			break
		}
	}
	c.Data["json"] = ""
	c.ServeJSON()
}

func (c *ServerController)Upload(){
	file, header, er := c.GetFile("file") // where <<this>> is the controller and <<file>> the id of your form field
    if er != nil {
        // get the filename
     	beego.Error(er)   
        
    }else{
    	fileName := header.Filename
    	// save to server
    	beego.Info(fileName)
    	beego.Info(file)
    	c.SaveToFile("file", "F://file.md")
    }
    c.Data["json"] = ""
	c.ServeJSON()
}