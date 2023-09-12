package events

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"sync"
)

type chanHolder struct {
	Ch interface{}
}

var (
	conn  *nats.Conn
	subCh []chanHolder
	mu    sync.Mutex
)

func Connect(server string) error {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil && len(subCh) > 0 {
		return nil
	}

	var err error
	conn, err = nats.Connect(server, nats.ReconnectBufSize(5*1024*1024))
	if err != nil {
		return err
	}
	subCh = make([]chanHolder, 0)
	return nil
}

func Subscribe[T any](topic string) (chan T, error) {
	ch := make(chan T)
	_, err := conn.Subscribe(topic, func(m *nats.Msg) {
		var msg T
		err := json.Unmarshal(m.Data, &msg)
		if err != nil {
			log.Printf("Error unmarshalling data: %v", err)
			return
		}
		ch <- msg
	})
	if err != nil {
		return nil, err
	}
	mu.Lock()
	subCh = append(subCh, chanHolder{Ch: ch})
	mu.Unlock()
	return ch, nil
}

func Publish[T any](msg T, topic string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return conn.Publish(topic, data)
}

func Close() {
	conn.Close()
	for _, holder := range subCh {
		if ch, ok := holder.Ch.(chan any); ok {
			close(ch)
		}
	}
}
