package src

import (
	"net/http"
	"code.google.com/p/go.net/websocket"
	"sync"
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
		websocket.Message.Send(ws,data)
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
	defer Close()
	for {
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			return
		}
	}
}
func xWebsocket(){
	http.Handle("/bifrost/demo/wsapi", websocket.Handler(CallBackServer))
}
