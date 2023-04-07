package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Subscriber struct {
	Addr  string
	Topic string
}

func (s *Subscriber) ReceiveMessages(topic string, number uint64, handler func(msg string)) error {
	conn, err := stomp.Dial("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", s.Addr, err)
	}
	sub, err := conn.Subscribe(topic, stomp.AckAuto,
		stomp.SubscribeOpt.Header("subscription-type", "MULTICAST"))
	if err != nil {
		return fmt.Errorf("cannot subscribe to %s: %v", topic, err)
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

func (s *Subscriber) ReceiveFrom(topic string, handler func(msg string)) error {
	return s.ReceiveMessages(topic, infinite, handler)
}

func (s *Subscriber) Receive(handler func(msg string)) error {
	if s.Topic == "" {
		return fmt.Errorf("no default topic specified")
	}
	return s.ReceiveFrom(s.Topic, handler)
}
