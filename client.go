package main 

import (
	"log"
	"time"
	"github.com/gorilla/websocket"
)

type wsClient struct  {
	socket *websocket.Conn
	exitFlag chan struct{}
}

func InitClient(attackUrl string) (c wsClient) {
	client, _, err := websocket.DefaultDialer.Dial(attackUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.socket = client
	c.exitFlag = make(chan struct{})
	return c
}

func (c wsClient) startListener(r chan WsMessage) {
	defer close(c.exitFlag)
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		var response WsMessage
		response.message = string(message)
		response.timeRecv = time.Now().Format("20060102150405")
		r <- response
	}
}

func (c wsClient) gracefullExit() {
	err := c.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
}

func (c wsClient) send(message string) {
	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write err:", err)
		return
	}
}

func (c wsClient) close() {
	c.socket.Close()
}

