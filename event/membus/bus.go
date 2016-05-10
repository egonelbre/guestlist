package membus

import (
	"sync"

	"github.com/egonelbre/guestlist/event"
)

type Bus struct {
	listeners map[int]func(event.Event)

	m    sync.RWMutex
	next int
}

func New() event.Bus {
	return &Bus{
		make(map[int]func(event.Event)),
		sync.RWMutex{},
		0,
	}
}

func (bus *Bus) Listen(fn func(event.Event)) (cancel func()) {
	bus.m.Lock()
	id := bus.next
	bus.next += 1
	bus.listeners[id] = fn
	bus.m.Unlock()

	return func() {
		bus.m.Lock()
		delete(bus.listeners, id)
		bus.m.Unlock()
	}
}

func (bus *Bus) Publish(e event.Event) {
	bus.m.RLock()
	for _, fn := range bus.listeners {
		fn(e)
	}
	bus.m.RUnlock()
}
