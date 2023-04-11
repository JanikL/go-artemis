package artemis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"math"
)

const infinite = math.MaxUint64

// Receiver receives values as gobs from a specified destination.
type Receiver struct {
	Addr string
	Dest string

	// PubSub configures the type of destination.
	// "true" for the publish-subscribe pattern (topics), "false" for the producer-consumer pattern (queues).
	PubSub bool
}

func (r *Receiver) ReceiveMessages(destination string, number uint64, handler func(msg any)) error {
	conn, err := stomp.Dial("tcp", r.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", r.Addr, err)
	}
	defer conn.Disconnect()
	sub, err := conn.Subscribe(destination, stomp.AckAuto, stomp.SubscribeOpt.Header("subscription-type", r.subType()))
	if err != nil {
		return fmt.Errorf("cannot receive from %s: %v", destination, err)
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

func (r *Receiver) ReceiveFrom(destination string, handler func(msg any)) error {
	return r.ReceiveMessages(destination, infinite, handler)
}

func (r *Receiver) Receive(handler func(msg any)) error {
	if r.Dest == "" {
		return fmt.Errorf("no default destination specified")
	}
	return r.ReceiveFrom(r.Dest, handler)
}

func (r *Receiver) subType() string {
	if r.PubSub {
		return "MULTICAST"
	} else {
		return "ANYCAST"
	}
}

func decode(message []byte) (any, error) {
	buff := bytes.NewBuffer(message)
	gob.Register(message)
	dec := gob.NewDecoder(buff)
	var msg any
	err := dec.Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	return msg, nil
}
