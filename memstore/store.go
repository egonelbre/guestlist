package memstore

import (
	"sync"
	"time"

	"github.com/egonelbre/event"
)

type eventlist struct {
	sync.Mutex
	items []event.Info
}

type Store struct {
	events    map[event.AggregateId]*eventlist
	publisher event.Publisher
	m         sync.Mutex
}

func New(pub event.Publisher) *Store {
	return &Store{
		make(map[event.AggregateId]*eventlist),
		pub,
		sync.Mutex{},
	}
}

func (store *Store) Save(id event.AggregateId, expectedVersion int64, events ...event.Event) error {
	store.m.Lock()
	list, ok := store.events[id]
	if !ok {
		list = &eventlist{}
		store.events[id] = list
	}
	store.m.Unlock()

	list.Lock()
	defer list.Unlock()

	version := int64(0)
	if len(list.items) > 0 {
		version = list.items[len(list.items)-1].Version
	}

	if version != expectedVersion && expectedVersion != -1 {
		return event.ConcurrencyError
	}

	for _, data := range events {
		version += 1
		event := event.Info{id, version, time.Now(), data}
		list.items = append(list.items, event)
		store.publisher.Publish(data)
	}

	return nil
}

func (store *Store) SaveChanges(changes event.Changes) error {
	return store.Save(changes.GetId(), changes.GetVersion(), changes.GetChanges()...)
}

func (store *Store) List(id event.AggregateId) (events []event.Info, found bool) {
	store.m.Lock()
	list, ok := store.events[id]
	store.m.Unlock()
	if !ok {
		return nil, false
	}

	return list.items, true
}
