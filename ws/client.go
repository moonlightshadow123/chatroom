package ws

import(
	"net/http"
	"encoding/json"
	"test_proj/chatRoom/utils"
	"test_proj/chatRoom/mymsg"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Name string			
}

func NewClient(w http.ResponseWriter, r *http.Request, name string)(*Client, error){
	conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
	if err != nil {
		utils.Error("Cannot setup WebSocket connection:", err)
		return nil, err
	}
	//beego.Info(models.User{Username:"nihao"})
	var client Client
	client.Name = name
	client.Conn = conn
	return &client, nil
}

func(client*Client)Close(){
	client.Conn.Close()
}


func(client*Client)OnFromWeb()(*mymsg.Message, error){
     _, msgData, err := client.Conn.ReadMessage()

    // 如果返回错误，就退出循环
    if err != nil{
        return nil, err
    }
    msg := Message{}
    if err := json.Unmarshal(msgData, &msg); err!=nil{
    	utils.Error("Json parse error!", err)
    	return nil, err
    }
    msg.Name = client.Name
	utils.Info("######WS-----------receive: ",msg)

	//如果没有错误，则把用户发送的信息放入message通道中
	/*
	msg := Message{Name: 	client.name,
					Type: 	0,
					Content:string(msgStr)}
	*/
	return &msg, nil
}


func(client*Client)OnFromManager(msg *mymsg.Message)(error){
	data, err := json.Marshal(msg)
	if err != nil {
		//beego.Error("Fail to marshal message:", err)
		return err
	}
	// fmt.Println("=======the json message is", string(data))	// 转换成字符串类型便于查看
	if err = client.Conn.WriteMessage(websocket.TextMessage, data); err!=nil {
		//beego.Error("Fail to write message")
		return err
	}
	return nil
}