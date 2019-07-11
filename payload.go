package main

import (
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Payload struct {
	Action string `json:"action"`
	Body string `json:"body"`
}

type Payloads struct {
    Payloads []Payload `json:"payloads"`
}

func readPayloadsFromFile(filePath string) (p Payloads) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
    	fmt.Println(err)
	}
	fmt.Println("Successfully Opened input file", filePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &p)
	return p
}