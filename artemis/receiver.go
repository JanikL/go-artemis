package artemis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
)

const infinite = math.MaxUint64

type Receiver interface {
	ReceiveMessages(destination string, number uint64, handler func(msg any)) error
	ReceiveFrom(destination string, handler func(msg any)) error
	Receive(handler func(msg any)) error
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
