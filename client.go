package main 

import (
	"log"
	"github.com/gorilla/websocket"
)

func InitClient(attackUrl string) (client *websocket.Conn) {
	client, _, err := websocket.DefaultDialer.Dial(attackUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return client
}

func setSocketListener(client *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := client.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

