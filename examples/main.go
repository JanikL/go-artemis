package main

import (
	"fmt"
	"github.com/JanikL/go-artemis/artemis"
	"log"
	"os"
	"time"
)

type TestDto struct {
	Message string
	Result  int
}

var brokerAddr = "localhost:61616"
var destination = "myDest"

func main() {
	if val, isPresent := os.LookupEnv("ARTEMIS_BROKER_ADDR"); isPresent {
		brokerAddr = val
	}

	log.Println("start receiver 1")
	go func() {
		var receiver1 artemis.Receiver = &artemis.Consumer{Addr: brokerAddr, Queue: destination}
		err := receiver1.Receive(handler1)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("start receiver 2")
	go func() {
		var receiver2 artemis.Receiver = &artemis.Consumer{Addr: brokerAddr, Queue: destination}
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
	var sender artemis.Sender = &artemis.Producer{Addr: brokerAddr, Queue: destination}
	for {
		err := sender.Send(TestDto{Message: "the answer is", Result: 42})
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func handler1(message interface{}) {
	if msg, ok := message.(TestDto); ok {
		fmt.Println("Receiver 1: ", msg.Message, msg.Result)
	} else {
		fmt.Println("Receiver 1: cannot assert type of message")
	}
}

func handler2(message interface{}) {
	if msg, ok := message.(TestDto); ok {
		fmt.Println("Receiver 2: ", msg.Message, msg.Result)
	} else {
		fmt.Println("Receiver 2: cannot assert type of message")
	}
}
