package main

import (
	"fmt"
	"github.com/JanikL/go-artemis/artemis"
	"log"
	"os"
	"time"
)

type TestDto struct {
	Message string `json:"message"`
	Result  int    `json:"result"`
}

var brokerAddr = "localhost:61616"
var destination = "myDest"

func main() {
	if val, isPresent := os.LookupEnv("ARTEMIS_BROKER_ADDR"); isPresent {
		brokerAddr = val
	}

	log.Println("start receiver 1")
	go func() {
		receiver1 := artemis.Receiver[TestDto]{
			Addr: brokerAddr,
			Dest: destination,
			Enc:  artemis.EncodingJson,
		}
		err := receiver1.Receive(handler1)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("start receiver 2")
	go func() {
		receiver2 := artemis.Receiver[TestDto]{
			Addr: brokerAddr,
			Dest: destination,
			Enc:  artemis.EncodingJson,
		}
		err := receiver2.Receive(handler2)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("start sending messages")
	err := sendInfiniteMessages()
	if err != nil {
		log.Fatal(err)
	}
}

func sendInfiniteMessages() error {
	sender := artemis.Sender{Addr: brokerAddr, Dest: destination, Enc: artemis.EncodingJson}
	for {
		err := sender.Send(TestDto{Message: "the answer is", Result: 42})
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func handler1(msg TestDto) {
	fmt.Println("receiver 1:", msg.Message, msg.Result)
}

func handler2(msg TestDto) {
	fmt.Println("receiver 2:", msg.Message, msg.Result)
}
