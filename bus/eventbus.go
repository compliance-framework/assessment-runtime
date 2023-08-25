package bus

import "sync"

type EventBusTunnel struct {
	subscribers []chan string
	lock        sync.RWMutex
}

func NewEventBusTunnel() *EventBusTunnel {
	// Listen to the event bus, probably with a channel, etc.
	// Assuming the events from the event bus flow through eventCh
	bus.Listen(eventCh)

	// Forward the event to all subscribers
	go func() {
		for event := range eventCh {
			et.lock.RLock()
			for _, ch := range et.subscribers {
				ch <- event
			}
			et.lock.RUnlock()
		}
	}()
	return et
}

// Other components of the runtime would use this to subscribe to events coming from the event bus
func (et *EventBusTunnel) Subscribe() chan string {
	et.lock.Lock()
	defer et.lock.Unlock()

	ch := make(chan string)
	et.subscribers = append(et.subscribers, ch)
	return ch
}
