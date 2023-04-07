package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Publisher struct {
	Addr  string
	Topic string
}

func (p *Publisher) SendTo(topic string, messages []string) error {
	conn, err := stomp.Dial("tcp", p.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", p.Addr, err)
	}
	for _, msg := range messages {
		err = conn.Send(topic, "text/plain", []byte(msg),
			stomp.SendOpt.Header("destination-type", "MULTICAST"))
		if err != nil {
			return fmt.Errorf("failed to send to %s: %v", topic, err)
		}
	}
	return conn.Disconnect()
}

func (p *Publisher) Send(messages []string) error {
	if p.Topic == "" {
		return fmt.Errorf("no default topic specified")
	}
	return p.SendTo(p.Topic, messages)
}

func (p *Publisher) SendMessage(msg string) error {
	return p.Send([]string{msg})
}
