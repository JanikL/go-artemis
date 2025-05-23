// Copyright (c) 2023 Janik Liebrecht
// Use of this source code is governed by the MIT License that can be found in the LICENSE file.

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

// A Sender sends messages to the artemis broker.
type Sender struct {

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

// SendTo sends messages to a specified destination.
func (s *Sender) SendTo(destination string, messages ...any) error {
	conn, err := stomp.Dial("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("could not connect to broker %s: %v", s.Addr, err)
	}
	defer conn.Disconnect()

	destType := s.destType()
	contentType := s.contentType()
	for _, msg := range messages {
		m, err := encode(msg, s.Enc)
		if err != nil {
			return fmt.Errorf("failed to encode message: %v: %v", msg, err)
		}
		opts := stomp.SendOpt.Header("destination-type", destType)
		err = conn.Send(destination, contentType, m, opts)
		if err != nil {
			return fmt.Errorf("could not send to destination %s: %v", destination, err)
		}
	}
	return nil
}

// Send sends messages to the default destination.
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

func (s *Sender) contentType() string {
	switch s.Enc {
	case EncodingGob:
		return "application/octet-stream"
	case EncodingJson:
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

func encode(message any, enc encoding) ([]byte, error) {
	switch enc {
	case EncodingGob:
		return encodeGob(message)
	case EncodingJson:
		return encodeJson(message)
	default:
		return nil, fmt.Errorf("unknown encoding: %v", enc)
	}
}

func encodeGob(message any) ([]byte, error) {
	gob.Register(message)
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&message) // Pass pointer to interface so Encode sees a value of interface type.
	if err != nil {
		return nil, fmt.Errorf("could not encode as gob: %v", err)
	}
	return buff.Bytes(), nil
}

func encodeJson(message any) ([]byte, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("could not marshal as json: %v", err)
	}
	return b, nil
}
