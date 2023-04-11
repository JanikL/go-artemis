package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Consumer struct {
	Addr  string
	Queue string
}

func (c *Consumer) ReceiveMessages(queue string, number uint64, handler func(msg any)) error {
	conn, err := stomp.Dial("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", c.Addr, err)
	}
	defer conn.Disconnect()
	sub, err := conn.Subscribe(queue, stomp.AckAuto,
		stomp.SubscribeOpt.Header("subscription-type", "ANYCAST"))
	if err != nil {
		return fmt.Errorf("cannot consume from %s: %v", queue, err)
	}
	var i uint64 = 0
	for ; number == infinite || i < number; i++ {
		msg := <-sub.C
		if msg.Err != nil {
			return fmt.Errorf("failed to receive a message: %v", msg.Err)
		}
		m, err := decode(msg.Body)
		if err != nil {
			return fmt.Errorf("failed to decode message: %v %v", msg.Header, err)
		}
		handler(m)
	}
	return nil
}

func (c *Consumer) ReceiveFrom(queue string, handler func(msg any)) error {
	return c.ReceiveMessages(queue, infinite, handler)
}

func (c *Consumer) Receive(handler func(msg any)) error {
	if c.Queue == "" {
		return fmt.Errorf("no default queue specified")
	}
	return c.ReceiveFrom(c.Queue, handler)
}
