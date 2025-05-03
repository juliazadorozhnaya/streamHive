package pubsub

import (
	"github.com/nats-io/nats.go"
	"log"
)

type Publisher struct {
	nc *nats.Conn
}

func NewNATS(url string) *Publisher {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	return &Publisher{nc: nc}
}

// Publish публикует сообщение в указанный subject
func (p *Publisher) Publish(subject string, payload []byte) error {
	return p.nc.Publish(subject, payload)
}

// Subscribe подписывается на subject и вызывает handler при получении данных
func (p *Publisher) Subscribe(subject string, handler func(data []byte)) error {
	_, err := p.nc.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	return err
}
