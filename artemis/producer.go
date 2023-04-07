package artemis

import (
	"fmt"
	"github.com/go-stomp/stomp/v3"
)

type Producer struct {
	Addr  string
	Queue string
}

func (p *Producer) SendTo(queue string, messages []string) error {
	conn, err := stomp.Dial("tcp", p.Addr)
	if err != nil {
		return fmt.Errorf("cannot connect to server %s: %v", p.Addr, err)
	}
	for _, msg := range messages {
		err = conn.Send(queue, "text/plain", []byte(msg),
			stomp.SendOpt.Header("destination-type", "ANYCAST"))
		if err != nil {
			return fmt.Errorf("failed to send to %s: %v", queue, err)
		}
	}
	return conn.Disconnect()
}

func (p *Producer) Send(messages []string) error {
	if p.Queue == "" {
		return fmt.Errorf("no default queue specified")
	}
	return p.SendTo(p.Queue, messages)
}

func (p *Producer) SendMessage(msg string) error {
	return p.Send([]string{msg})
}
