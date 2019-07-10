package main 

import (
	"log"
	"os"
	"os/signal"
	"fmt"
	"time"
	"strconv"
	"github.com/gorilla/websocket"
)

func InitClient(attackUrl string) (client *websocket.Conn) {
	client, _, err := websocket.DefaultDialer.Dial(attackUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return client
}

func handleMessage(client *websocket.Conn, done chan struct{}) {
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

func Fuzz(client *websocket.Conn, attackUrl string, payload Payloads,  done chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for i := 0; i < len(payload.Payloads); i++ {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			break;
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		default:
			PeformAction(client, payload.Payloads[i])
		}
	}
}

func PeformAction(client *websocket.Conn, payload Payload) {
	fmt.Println("handeling action : ", payload.Action);
	switch action := payload.Action; action {
	case "send":
		err := client.WriteMessage(websocket.TextMessage, []byte(payload.Body))
		if err != nil {
			log.Println("write:", err)
			return
		}	
		fmt.Println("sent ", payload.Body);
	case "wait":
		ms, err := strconv.Atoi(payload.Body)
		if err != nil {
		    log.Println("delay:", err)
			return
		}
		fmt.Println("wait for ", payload.Body);
		time.Sleep(time.Duration(ms) * time.Second)
	}
}

func Run(attackUrl, inputPath string) {
	var payloads = readPayloadsFromFile(inputPath)
	var client = InitClient(attackUrl)
	defer client.Close()
	done := make(chan struct{})
	// set up event handler
	go handleMessage(client, done)
	fmt.Println("created client, starting fuzz on ", attackUrl);
	Fuzz(client, attackUrl, payloads, done)
}
