package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Publisher struct {
	Addr  string
	Topic string
}

func (p *Publisher) SendTo(topic string, messages []any) error {
	conn, err := stomp.Dial("tcp", p.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", p.Addr, err)
	}
	for _, msg := range messages {
		m, err := encode(msg)
		if err != nil {
			return fmt.Errorf("failed to encode message %v: %v", msg, err)
		}
		err = conn.Send(topic, "text/plain", m,
			stomp.SendOpt.Header("destination-type", "MULTICAST"))
		if err != nil {
			return fmt.Errorf("failed to send to %s: %v", topic, err)
		}
	}
	return conn.Disconnect()
}

func (p *Publisher) Send(messages []any) error {
	if p.Topic == "" {
		return fmt.Errorf("no default topic specified")
	}
	return p.SendTo(p.Topic, messages)
}

func (p *Publisher) SendMessage(msg any) error {
	return p.Send([]any{msg})
}
