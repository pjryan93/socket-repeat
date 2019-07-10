package main

import (
	"os"
	"os/signal"
	"fmt"
	"time"
	"strconv"
	"log"
	"github.com/gorilla/websocket"
)


type Fuzzer struct {
	socket  *websocket.Conn
	payloads Payloads
	finishedFlag chan struct{}
}

func CreateFuzzer(url, inputPath string) (socketFuzzer Fuzzer){
	var payloads = readPayloadsFromFile(inputPath)
	var client = InitClient(url)
	done := make(chan struct{})
	socketFuzzer.socket = client
	socketFuzzer.payloads = payloads
	socketFuzzer.finishedFlag = done
	go setSocketListener(socketFuzzer.socket, socketFuzzer.finishedFlag)
	return socketFuzzer
}


func Fuzz(fuzzer Fuzzer) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for i := 0; i < len(fuzzer.payloads.Payloads); i++ {
		select {
		case <-fuzzer.finishedFlag:
			return
		case <-interrupt:
			log.Println("interrupt")
			break;
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := fuzzer.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-fuzzer.finishedFlag:
			case <-time.After(time.Second):
			}
			return
		default:
			peformAction(fuzzer.socket, fuzzer.payloads.Payloads[i])
		}
	}
}

func peformAction(client *websocket.Conn, payload Payload) {
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
	Fuzzer := CreateFuzzer(attackUrl, inputPath)
	defer Fuzzer.socket.Close()
	fmt.Println("created client, starting fuzz on ", attackUrl);
	Fuzz(Fuzzer)
}
