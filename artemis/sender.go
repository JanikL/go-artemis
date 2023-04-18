package artemis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type encoding int

const (
	EncodingGob encoding = iota
	EncodingJson
)

// Sender sends values as gobs to a specified destination.
type Sender struct {
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

func (s *Sender) SendTo(destination string, messages ...any) error {
	conn, err := stomp.Dial("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", s.Addr, err)
	}
	defer conn.Disconnect()
	destType := s.destType()
	for _, msg := range messages {
		m, err := encode(msg, s.Enc)
		if err != nil {
			return fmt.Errorf("failed to encode message: %v: %v", msg, err)
		}
		err = conn.Send(destination, "text/plain", m,
			stomp.SendOpt.Header("destination-type", destType))
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

func encode(message any, enc encoding) ([]byte, error) {
	switch enc {
	case EncodingGob:
		return encodeGob(message)
	case EncodingJson:
		return encodeJson(message)
	default:
		panic(fmt.Sprint("unknown encoding", enc))
	}
}

func encodeGob(message any) ([]byte, error) {
	gob.Register(message)
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&message) // Pass pointer to interface so Encode sees a value of interface type.
	if err != nil {
		return nil, fmt.Errorf("encode error: %v", err)
	}
	return buff.Bytes(), nil
}

func encodeJson(message any) ([]byte, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode error: %v", err)
	}
	return b, nil
}
