package src

import (
	"net/http"
	"golang.org/x/net/websocket"
	"sync"
	"log"
)

var websocketList map[int]*websocket.Conn
var LastWebsocketID int = 0
var l sync.RWMutex

func init()  {
	websocketList = make(map[int]*websocket.Conn,0)
	xWebsocket()
}

func sendMsgToWsClient(data interface{}){
	l.RLock()
	defer l.RUnlock()
	for _,ws := range websocketList{
		err := websocket.Message.Send(ws,data)
		if err != nil{
			log.Println("sendMsgToWsClient err:",err)
		}
	}
}

func getWebSocketID() int{
	l.Lock()
	ID := LastWebsocketID
	LastWebsocketID++
	l.Unlock()
	return ID
}

func CallBackServer(ws *websocket.Conn) {
	var msg string
	var err error
	ID := getWebSocketID()
	l.Lock()
	websocketList[ID] = ws
	l.Unlock()
	var Close = func() {
		ws.Close()
		l.Lock()
		delete(websocketList, ID)
		l.Unlock()
	}
	websocket.Message.Send(ws,"connect success")
	defer Close()
	for {
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			log.Println("ws CallBackServer err:",err)
			return
		}
	}
}
func xWebsocket(){
	http.Handle("/bifrost/demo/wsapi", websocket.Handler(CallBackServer))
}
