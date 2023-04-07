package artemis

import "math"

const infinite = math.MaxUint64

type Receiver interface {
	ReceiveMessages(destination string, number uint64, handler func(msg string)) error
	ReceiveFrom(destination string, handler func(msg string)) error
	Receive(handler func(msg string)) error
}
