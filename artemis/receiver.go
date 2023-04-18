package artemis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"math"
)

const infinite = math.MaxUint64

// Receiver receives values as gobs from a specified destination.
type Receiver[T any] struct {
	Addr string
	Dest string

	// PubSub configures the type of destination.
	// "true" for the publish-subscribe pattern (topics),
	// "false" for the producer-consumer pattern (queues).
	PubSub bool

	// Enc specifies the encoding.
	// The default encoding is gob.
	Enc encoding
}

func (r *Receiver[T]) ReceiveMessages(destination string, number uint64, handler func(msg T)) error {
	conn, err := stomp.Dial("tcp", r.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", r.Addr, err)
	}
	defer conn.Disconnect()
	sub, err := conn.Subscribe(destination, stomp.AckAuto,
		stomp.SubscribeOpt.Header("subscription-type", r.subType()))
	if err != nil {
		return fmt.Errorf("cannot receive from %s: %v", destination, err)
	}
	var i uint64 = 0
	for ; number == infinite || i < number; i++ {
		msg := <-sub.C
		if msg.Err != nil {
			return fmt.Errorf("failed to receive a message: %v", msg.Err)
		}
		m, err := decode[T](msg.Body, r.Enc)
		if err != nil {
			return fmt.Errorf("failed to decode message: %v %v", msg.Header, err)
		}
		handler(m)
	}
	return nil
}

func (r *Receiver[T]) ReceiveFrom(destination string, handler func(msg T)) error {
	return r.ReceiveMessages(destination, infinite, handler)
}

func (r *Receiver[T]) Receive(handler func(msg T)) error {
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

func decode[T any](message []byte, enc encoding) (T, error) {
	switch enc {
	case EncodingGob:
		return decodeGob[T](message)
	case EncodingJson:
		return decodeJson[T](message)
	default:
		panic(fmt.Sprint("unknown encoding", enc))
	}
}

func decodeGob[T any](message []byte) (T, error) {
	gob.Register(*new(T))
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var msg any
	err := dec.Decode(&msg)
	if err != nil {
		return *new(T), fmt.Errorf("decode error: %v", err)
	}
	m := msg.(T)
	return m, nil
}

func decodeJson[T any](message []byte) (T, error) {
	var msg T
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return msg, fmt.Errorf("decode error: %v", err)
	}
	return msg, nil
}
