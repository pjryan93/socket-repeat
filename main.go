package main 

import (
	"flag"
	"fmt"
)

func main() {
	var attackUrl string
	var fileName string
	flag.StringVar(&attackUrl, "url", "ws://localhost:8080", "url to open websocket")
	flag.StringVar(&fileName, "f", "payloads.json", "json of payloads")
	flag.Parse()
	fmt.Println("Initing socket to ", attackUrl);
	Run(attackUrl, fileName)
}