package bus

import (
	"encoding/json"
	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	conn           *nats.Conn
	configChannels map[string][]chan<- config.Config
	mu             sync.Mutex
)

func Connect(server string) error {
	var err error
	conn, err = nats.Connect(server)
	if err != nil {
		return err
	}
	configChannels = make(map[string][]chan<- config.Config)
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
	configChannels["configuration"] = append(configChannels["configuration"], ch)
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
	// Close the connection
	conn.Close()
}
