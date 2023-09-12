package events

import (
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Message struct {
	Text string `json:"text"`
}

func TestBus(t *testing.T) {
	s := natsserver.RunServer(&natsserver.DefaultTestOptions)
	defer s.Shutdown()

	err := Connect(nats.DefaultURL)
	assert.NoError(t, err)

	topic := "test"
	msg := Message{Text: "Hello, World!"}

	ch, err := Subscribe[Message](topic)
	assert.NoError(t, err)
	assert.NotNil(t, ch)

	err = Publish(msg, topic)
	assert.NoError(t, err)

	received := <-ch
	assert.Equal(t, msg.Text, received.Text)

	Close()
}
