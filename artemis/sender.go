package artemis

type Sender interface {
	SendTo(destination string, messages []string) error
	Send(messages []string) error
	SendMessage(msg string) error
}
