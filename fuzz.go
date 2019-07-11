package main

import (
	"os"
	"os/signal"
	"fmt"
	"time"
	"strconv"
	"log"
)

type Fuzzer struct {
	socket  wsClient
	payloads Payloads
	responses chan WsMessage
	sent chan WsMessage
}

func createFuzzer(url, inputPath string) (socketFuzzer Fuzzer){
	var payloads = readPayloadsFromFile(inputPath)
	var client = InitClient(url)
	socketFuzzer.socket = client
	socketFuzzer.payloads = payloads
	socketFuzzer.responses = make(chan WsMessage)
	socketFuzzer.sent = make(chan WsMessage)
	return socketFuzzer
}

func (fuzzer Fuzzer) fuzz() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for i := 0; i < len(fuzzer.payloads.Payloads); i++ {
		select {
		case <-fuzzer.socket.exitFlag:
			return
		case <-interrupt:
			log.Println("interrupt")
			fuzzer.socket.gracefullExit()
			select {
			case <-fuzzer.socket.exitFlag:
			case <-time.After(time.Second):
			}
			return
		default:
			fuzzer.handlePayload(fuzzer.payloads.Payloads[i])
		}
	}
}

func (fuzzer Fuzzer) handlePayload(payload Payload) {
	fmt.Println("handeling payload action : ", payload.Action);
	switch action := payload.Action; action {
	case "send":
		fuzzer.socket.send(payload.Body)
		fuzzer.sent <- WsMessage{payload.Body, time.Now().Format("20060102150405")}
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
	fuzzer := createFuzzer(attackUrl, inputPath)
	go fuzzer.socket.startListener(fuzzer.responses)
	defer fuzzer.socket.close()
	fmt.Println("created client, started listening, begining fuzz on ", attackUrl);
	fuzzer.fuzz()
	fmt.Println(<-fuzzer.responses)
}
