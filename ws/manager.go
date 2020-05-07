package ws

import(
	"fmt"
	"test_proj/chatRoom/mymsg"
	"test_proj/chatRoom/models"
	"test_proj/chatRoom/utils"
	//"github.com/gorilla/websocket"
)

type Message = mymsg.Message
var(
	manInst *Manager
)

type Manager struct{
	Join 	chan *Client
	Leave	chan *Client
	Msgs chan *Message
	Clients map [string] *Client
}

func GetInst()(man*Manager){
	if manInst == nil{
		manInst = &Manager{
			Join: make(chan *Client, 10),
			Leave: make(chan *Client, 10),
			Msgs : make(chan *Message, 10),
			Clients: make(map[string]*Client),
		}
	}
	return manInst
}

func(man*Manager)Run(){
	go man.broadcast()
}

func(man*Manager)broadcast(){
	for {
		// 哪个case可以执行，则转入到该case。都不可执行，则堵塞。
		select {
			// 消息通道中有消息则执行，否则堵塞
			case msg := <-man.Msgs:
				man.onMsg(msg)
			// 有用户加入
			case client := <-man.Join:
				man.onJoin(client)
			// 有用户退出
			case client := <-man.Leave:
				man.onLeave(client)
		}
	}
}

func (man*Manager)TypeOper(msg * Message, dbmsg *models.Message)(error){
	if msg.Type == 3{
		num, err := models.RecallMsg(dbmsg)
		if num != 1 || err != nil{
			return fmt.Errorf("Recall msg error. dbmsg:%s  num:%d err:%s", dbmsg, num, err)
		}
	}
	return nil
}

func (man*Manager)SyncMsgToDB(msg * Message)(error){
	dbmsg,err := mymsg.MsgToDbMsg(msg)
	if err != nil{
		return err
	}
	err1:= man.TypeOper(msg, dbmsg)
	if err1 != nil{
		return err1
	}
	err2 := models.NewMsg(dbmsg)
	if err2 != nil{
		return err2
	}
	fullmsg := mymsg.DbMsgToMsg(dbmsg) // fill id, date 
	mymsg.Copy(fullmsg, msg)
	return nil
}

func(man*Manager)onJoin(client *Client){
	str := fmt.Sprintf("broadcaster-----------%s join in the chat room\n", client.Name)
	utils.Info(str)
	
	man.Clients[client.Name] = client	// 将用户加入映射

	// 将用户加入消息放入消息通道
	content := fmt.Sprintf("%s join in, there are %d preson in room", client.Name, len(man.Clients))
	msg := Message{	Name:	client.Name, 
					Type:	1, 
					Content:content}
	// 此处要设置有缓冲的通道。因为这是goroutine自己从通道中发送并接受数据。
	// 若是无缓冲的通道，该goroutine发送数据到通道后就被锁定，需要数据被接受后才能解锁，而恰恰接受数据的又只能是它自己
	man.Msgs <- &msg
}

func(man*Manager)onLeave(client *Client){
	str := fmt.Sprintf("broadcaster-----------%s leave the chat room\n", client.Name)
	utils.Info(str)

	// 如果该用户已经被删除
	_, has := man.Clients[client.Name]
	if !has {
		utils.Info("the client had leaved, client's name:"+client.Name)
		return
		// break
	}
	delete(man.Clients, client.Name)	// 将用户从映射中删除

	// 将用户退出消息放入消息通道
	var msg Message
	msg.Name = client.Name
	msg.Type = 2
	msg.Content = fmt.Sprintf("%s leave, there are %d preson in room", client.Name, len(man.Clients))
	man.Msgs <- &msg
}

func(man*Manager)onMsg(msg *Message){
	// str := fmt.Sprintf("broadcaster-----------%s send message: %s\n", msg.Name, msg.Content)
	// beego.Info(str)

	if err := man.SyncMsgToDB(msg); err != nil{
		utils.Error("Fail to syn MSG:'", msg ,"' to db.", err)
		return
	}
	str := fmt.Sprintf("Msg:'%s' synced to DB. Now broadcasting it!", msg)
	utils.Info(str)
	// 将某个用户发出的消息发送给所有用户
	for _, client := range man.Clients {
		// 将数据编码成json形式，data是[]byte类型
		// json.Marshal()只会编码结构体中公开的属性(即大写字母开头的属性)
		client.OnFromManager(msg)
	}
}

func(man*Manager)OnUpload(name string, ori_filename string, file_path string, addr string){
	var msg Message
	msg.Name = name
	msg.Type = 4
	msg.Content = fmt.Sprintf("%s uploaded a file.", name)
	msg.Related = ori_filename + ";"+ file_path + ";"+addr
	man.Msgs <- &msg
}