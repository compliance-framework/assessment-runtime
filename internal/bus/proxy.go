package bus

import (
	"encoding/json"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	conn           *nats.Conn
	configChannels []chan<- config.Config
	mu             sync.Mutex
)

func Connect(server string) error {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil && configChannels != nil {
		return nil
	}

	var err error
	conn, err = nats.Connect(server, nats.ReconnectBufSize(5*1024*1024))
	if err != nil {
		return err
	}
	configChannels = make([]chan<- config.Config, 0)
	return nil
}

func SubToConfig(ch chan<- config.Config) error {
	_, err := conn.Subscribe("configuration", func(m *nats.Msg) {
		var cfg config.Config
		err := json.Unmarshal(m.Data, &cfg)
		if err != nil {
			log.Printf("Error unmarshalling data: %v", err)
			return
		}
		ch <- cfg
	})
	if err != nil {
		return err
	}
	mu.Lock()
	configChannels = append(configChannels, ch)
	mu.Unlock()
	return nil
}

func PubConfig(cfg config.Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return conn.Publish("configuration", data)
}

func Close() {
	conn.Close()
}
