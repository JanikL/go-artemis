package artemis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

// Sender sends values as gobs to a specified destination.
type Sender struct {
	Addr string
	Dest string

	// PubSub configures the type of destination.
	// "true" for the publish-subscribe pattern (topics), "false" for the producer-consumer pattern (queues).
	PubSub bool
}

func (s *Sender) SendTo(destination string, messages ...any) error {
	conn, err := stomp.Dial("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", s.Addr, err)
	}
	defer conn.Disconnect()
	destType := s.destType()
	for _, msg := range messages {
		m, err := encode(msg)
		if err != nil {
			return fmt.Errorf("failed to encode message: %v: %v", msg, err)
		}
		err = conn.Send(destination, "text/plain", m, stomp.SendOpt.Header("destination-type", destType))
		if err != nil {
			return fmt.Errorf("failed to send to %s: %v", destination, err)
		}
	}
	return nil
}

func (s *Sender) Send(messages ...any) error {
	if s.Dest == "" {
		return fmt.Errorf("no default destination specified")
	}
	return s.SendTo(s.Dest, messages...)
}

func (s *Sender) destType() string {
	if s.PubSub {
		return "MULTICAST"
	} else {
		return "ANYCAST"
	}
}

func encode(message any) ([]byte, error) {
	buff := bytes.Buffer{}
	gob.Register(message)
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&message) // Pass pointer to interface so Encode sees a value of interface type.
	if err != nil {
		return nil, fmt.Errorf("encode error: %v", err)
	}
	return buff.Bytes(), nil
}
