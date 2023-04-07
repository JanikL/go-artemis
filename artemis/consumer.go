package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Consumer struct {
	Addr  string
	Queue string
}

func (c *Consumer) ReceiveMessages(queue string, number uint64, handler func(msg string)) error {
	conn, err := stomp.Dial("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", c.Addr, err)
	}
	sub, err := conn.Subscribe(queue, stomp.AckAuto,
		stomp.SubscribeOpt.Header("subscription-type", "ANYCAST"))
	if err != nil {
		return fmt.Errorf("cannot consume from %s: %v", queue, err)
	}
	var i uint64 = 0
	for ; number == infinite || i < number; i++ {
		msg := <-sub.C
		if msg.Err != nil {
			return msg.Err
		}
		handler(string(msg.Body))
	}
	return conn.Disconnect()
}

func (c *Consumer) ReceiveFrom(queue string, handler func(msg string)) error {
	return c.ReceiveMessages(queue, infinite, handler)
}

func (c *Consumer) Receive(handler func(msg string)) error {
	if c.Queue == "" {
		return fmt.Errorf("no default queue specified")
	}
	return c.ReceiveFrom(c.Queue, handler)
}
