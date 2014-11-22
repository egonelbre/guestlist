package invitation

import (
	"encoding/gob"

	"github.com/egonelbre/event"
)

type Created struct {
	Id   event.AggregateId
	Name string
}

type Accepted struct {
	Id event.AggregateId
}

type Declined struct {
	Id event.AggregateId
}

func init() {
	gob.Register(Created{})
	gob.Register(Accepted{})
	gob.Register(Declined{})
}
