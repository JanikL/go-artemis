package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Producer struct {
	Addr  string
	Queue string
}

func (p *Producer) SendTo(queue string, messages ...any) error {
	conn, err := stomp.Dial("tcp", p.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", p.Addr, err)
	}
	defer conn.Disconnect()
	for _, msg := range messages {
		m, err := encode(msg)
		if err != nil {
			return fmt.Errorf("failed to encode message: %v: %v", msg, err)
		}
		err = conn.Send(queue, "text/plain", m,
			stomp.SendOpt.Header("destination-type", "ANYCAST"))
		if err != nil {
			return fmt.Errorf("failed to send to %s: %v", queue, err)
		}
	}
	return nil
}

func (p *Producer) Send(messages ...any) error {
	if p.Queue == "" {
		return fmt.Errorf("no default queue specified")
	}
	return p.SendTo(p.Queue, messages...)
}

func (p *Producer) SendMessage(msg any) error {
	return p.Send(msg)
}
