package main 

import (
	"log"
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

func (c wsClient) startListener() {
	defer close(c.exitFlag)
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
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

