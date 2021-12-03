package nsqstore

import "github.com/nsqio/go-nsq"

type Nsq struct {
	producer *nsq.Producer
}

type NsqI interface {
	Publish(topic string, payload []byte) error
}

func Init(p *nsq.Producer) NsqI {
	return &Nsq{
		producer: p,
	}
}

// Отправить сообщение в очередь
func (n *Nsq) Publish(topic string, payload []byte) error {
	return n.producer.Publish(topic, payload)
}
