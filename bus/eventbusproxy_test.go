package bus

import (
	"github.com/nats-io/nats.go"
	"testing"
)

func TestNew(t *testing.T) {
	ebp, err := NewEventBusProxy(nats.DefaultURL)
	if ebp == nil || err != nil {
		t.Errorf("New() failed, expected non-nil EventBusProxy, got nil, error: %v", err)
	}
}

func TestSubscribe(t *testing.T) {
	ebp, _ := NewEventBusProxy(nats.DefaultURL)
	ch := make(chan interface{})
	err := ebp.Subscribe("test.subject", ch)
	if err != nil {
		t.Errorf("Subscribe() failed, expected no error, got error: %v", err)
	}
}

func TestPublish(t *testing.T) {
	ebp, _ := NewEventBusProxy(nats.DefaultURL)
	err := ebp.Publish("test.subject", []byte("test data"))
	if err != nil {
		t.Errorf("Publish() failed, expected no error, got error: %v", err)
	}
}

func TestClose(t *testing.T) {
	ebp, _ := NewEventBusProxy(nats.DefaultURL)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Close() failed, expected no panic, got panic: %v", r)
		}
	}()
	ebp.Close()
}

func TestSubscribeAndPublish(t *testing.T) {
	ebp, _ := NewEventBusProxy(nats.DefaultURL)
	ch := make(chan interface{})
	subject := "test.subject"
	data := []byte("test data")

	err := ebp.Subscribe(subject, ch)
	if err != nil {
		t.Errorf("Subscribe() failed, expected no error, got error: %v", err)
	}

	go func() {
		err := ebp.Publish(subject, data)
		if err != nil {
			t.Errorf("Publish() failed, expected no error, got error: %v", err)
		}
	}()

	receivedData := <-ch
	if string(receivedData.([]byte)) != string(data) {
		t.Errorf("SubscribeAndPublish() failed, expected %v, got %v", string(data), string(receivedData.([]byte)))
	}
}
