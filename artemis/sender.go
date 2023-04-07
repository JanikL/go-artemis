package artemis

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Sender interface {
	SendTo(destination string, messages []any) error
	Send(messages []any) error
	SendMessage(msg any) error
}

func encode(message any) ([]byte, error) {
	buff := bytes.Buffer{}
	gob.Register(message)
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(&message) // TODO why pointer
	if err != nil {
		return nil, fmt.Errorf("encode error: %v", err)
	}
	return buff.Bytes(), nil
}
