package artemis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

const unlimited = 0

// A Receiver receives messages of type T from the artemis broker.
type Receiver[T any] struct {

	// Addr is the address of the broker.
	Addr string

	// Dest is the default destination.
	Dest string

	// PubSub configures the type of destination.
	// 'true' for the publish-subscribe pattern (topics),
	// 'false' for the producer-consumer pattern (queues).
	PubSub bool

	// Enc specifies the encoding.
	// The default encoding is gob.
	Enc encoding
}

// ReceiveMessages receives a specified number of messages from a specified destination.
// If number is set to 0, it will receive an unlimited number of messages.
func (r *Receiver[T]) ReceiveMessages(destination string, number uint64, handler func(msg *T)) error {
	conn, err := stomp.Dial("tcp", r.Addr)
	if err != nil {
		return fmt.Errorf("could not connect to broker %s: %v", r.Addr, err)
	}
	defer conn.Disconnect()
	sub, err := conn.Subscribe(destination, stomp.AckAuto,
		stomp.SubscribeOpt.Header("subscription-type", r.subType()))
	if err != nil {
		return fmt.Errorf("could not subscribe to queue %s: %v", destination, err)
	}
	var i uint64 = 0
	for ; number == unlimited || i < number; i++ {
		msg := <-sub.C
		if msg.Err != nil {
			return fmt.Errorf("received an error: %v", msg.Err)
		}
		m, err := decode[T](msg.Body, r.Enc)
		if err != nil {
			return fmt.Errorf("failed to decode message: %v: %v", msg.Header, err)
		}
		handler(m)
	}
	return nil
}

// ReceiveFrom receives messages from a specified destination.
func (r *Receiver[T]) ReceiveFrom(destination string, handler func(msg *T)) error {
	return r.ReceiveMessages(destination, unlimited, handler)
}

// Receive receives messages from the default destination.
func (r *Receiver[T]) Receive(handler func(msg *T)) error {
	if r.Dest == "" {
		return fmt.Errorf("no default destination specified")
	}
	return r.ReceiveFrom(r.Dest, handler)
}

func (r *Receiver[T]) subType() string {
	if r.PubSub {
		return "MULTICAST"
	} else {
		return "ANYCAST"
	}
}

func decode[T any](message []byte, enc encoding) (*T, error) {
	switch enc {
	case EncodingGob:
		return decodeGob[T](message)
	case EncodingJson:
		return decodeJson[T](message)
	default:
		return nil, fmt.Errorf("unknown encoding: %v", enc)
	}
}

func decodeGob[T any](message []byte) (*T, error) {
	gob.Register(*new(T))
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var msg any
	err := dec.Decode(&msg)
	if err != nil {
		return nil, fmt.Errorf("could not decode gob: %v", err)
	}
	m := msg.(T)
	return &m, nil
}

func decodeJson[T any](message []byte) (*T, error) {
	var msg T
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}
	return &msg, nil
}
