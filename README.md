# go-artemis
Go language implementation of a simple ActiveMQ Artemis client library. The library offers functions 
for sending and receiving messages to/from the Artemis Broker using STOMP.

Features:
- Send and receive messages to/from queues and topics
- Encode messages as JSON or Gob; other encodings can be added easily

Limitations:
- A `Receiver` object can only receive messages of one type

## Usage
The import path for the package is `github.com/JanikL/go-artemis/artemis`.

### How to send and receive messages
Below is a simple example of how to send and receive messages of type `string` to/from a queue.
```go
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
```
### How to use topics
Queues are used by default. To use topics, set the `PubSub` field of the `Sender` and `Receiver` 
objects to `true`.

### How to use JSON encoding
Gob encoding is used by default. To encode messages as JSON, set the `Enc` field of the `Sender` 
and `Receiver` objects to `artemis.EncodingJson`.

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please feel 
free to open an issue or submit a pull request on the GitHub repository.

## License
Copyright (c) 2023 Janik Liebrecht

This project is licensed under the terms of the MIT License. See the [LICENSE](LICENSE.txt) file for 
details.
