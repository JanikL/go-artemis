package main

import (
	"fmt"
	"github.com/JanikL/go-artemis/artemis"
	"time"
)

const (
	brokerAddr  = "localhost:61616"
	destination = "myQueue"
)

func main() {

	// create a receiver that receives messages of type string
	receiver := artemis.Receiver[string]{Addr: brokerAddr, Dest: destination}

	// create a message handler
	handler := func(msg *string) {
		fmt.Println(*msg)
	}

	// start receiving messages in a separate goroutine
	go func() {
		err := receiver.Receive(handler)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// create a sender and send two messages
	sender := artemis.Sender{Addr: brokerAddr, Dest: destination}
	err := sender.Send("Hello", "Artemis")
	if err != nil {
		fmt.Println(err)
	}

	// wait a second to ensure the messages are received before the program exits
	time.Sleep(1 * time.Second)
}
