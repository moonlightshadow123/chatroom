package mymsg

import(
	"time"
	//"encoding/json"
	"test_proj/chatRoom/models"
	"test_proj/chatRoom/utils"
)

var(
	form = "2006-01-02 15:04:05"
	location = "Asia/Chongqing"
)

// 0表示用户发布消息；
// 1表示用户进入；
// 2表示用户退出; 
// 3表示用户撤回消息; 
// 4表示开始文件传输；
// 5表示删除文件
// 6表示发起语音聊天
// 7表示接受语音聊天
// 8表示结束语音聊天
// 9表示发起视频聊天
// 10表示接受视频聊天
// 11表示结束视频聊天
type Message struct {
	Type byte		`json:"type"`		
	Name string		`json:"name"`		// 用户名称
	Content string	`json:"message"`	// 消息
	Id int64 		`json:"id"`
	Related  string `json:"related"`
	Date string 	`json:"date"`
}


func DbMsgToMsg(dbmsg *models.Message)(*Message){
	msg := Message{ Id: 		dbmsg.Id,
					Type: 		byte(dbmsg.Type),
					Content:	dbmsg.Content,
					Related:	dbmsg.Related,
					Name:		dbmsg.User.Username,
					Date:		dbmsg.Date.Format(form)}
	return &msg
}

func MsgToDbMsg(msg*Message)(*models.Message, error){
	var date time.Time
	var err error
	if msg.Date == ""{
		date = time.Now()
		loc, e:= time.LoadLocation(location)
		if e !=nil{
			utils.Error("time.LoadLocation error, loc = =", location)
		}else{
			date = date.In(loc)
		}
	}else{
		date,err = time.Parse(form, msg.Date)
		if err != nil{
			return nil, err
		}
	}
	dbmsg := models.Message{Type: 		int(msg.Type),
							Content:	msg.Content,
							Related:	msg.Related,
							Date:		date,
							User:		&models.User{Username:msg.Name}}
	return &dbmsg, nil
}

func Copy(srcmsg*Message, trgmsg*Message){
	trgmsg.Type = srcmsg.Type
	trgmsg.Name = srcmsg.Name
	trgmsg.Content = srcmsg.Content
	trgmsg.Related = srcmsg.Related
	trgmsg.Id = srcmsg.Id
	trgmsg.Date = srcmsg.Date
}