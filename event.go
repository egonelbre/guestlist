package event

import (
	"errors"
	"time"
)

type Info struct {
	Id      AggregateId
	Version int64
	Date    time.Time
	Data    Event
}

type Event interface{}

type Store interface {
	Save(id AggregateId, expectedVersion int64, events ...Event) error
	SaveChanges(changes Changes) error
	List(id AggregateId) ([]Info, bool)
}

type Changes interface {
	GetId() AggregateId
	GetVersion() int64
	GetChanges() []Event
}

type Publisher interface {
	Publish(event Event)
}

type Bus interface {
	Listen(fn func(Event)) (cancel func())
	Publisher
}

var ConcurrencyError = errors.New("concurrency error")
