package artemis

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Sender interface {
	SendTo(destination string, messages ...any) error
	Send(messages ...any) error
}

func encode(message any) ([]byte, error) {
	buff := bytes.Buffer{}
	gob.Register(message)
	enc := gob.NewEncoder(&buff)

	// Pass pointer to interface so Encode sees a value of interface type.
	err := enc.Encode(&message)
	if err != nil {
		return nil, fmt.Errorf("encode error: %v", err)
	}
	return buff.Bytes(), nil
}
