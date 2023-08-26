package bus

import (
	"github.com/nats-io/nats.go"
)

type EventBusProxy struct {
	conn     *nats.Conn
	channels map[string]chan<- interface{}
}

func NewEventBusProxy(server string) (*EventBusProxy, error) {
	conn, err := nats.Connect(server)
	if err != nil {
		return nil, err
	}
	return &EventBusProxy{
		conn:     conn,
		channels: make(map[string]chan<- interface{}),
	}, nil
}

func (ebp *EventBusProxy) Subscribe(subject string, ch chan<- interface{}) error {
	_, err := ebp.conn.Subscribe(subject, func(m *nats.Msg) {
		ch <- m.Data
	})
	if err != nil {
		return err
	}
	ebp.channels[subject] = ch
	return nil
}

func (ebp *EventBusProxy) Publish(subject string, data interface{}) error {
	return ebp.conn.Publish(subject, data.([]byte))
}

func (ebp *EventBusProxy) Close() {
	ebp.conn.Close()
}
