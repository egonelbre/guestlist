package invitation

import (
	"fmt"

	"github.com/egonelbre/event"
)

type Aggregate struct {
	event.Aggregate
	Info
}

type Info struct {
	Name     string
	Accepted bool
	Declined bool
}

func (invite *Aggregate) Init(name string) {
	invite.Apply(Created{invite.Id, name})
}

func (invite *Aggregate) Accept() error {
	if invite.Declined {
		return fmt.Errorf("%s already declined", invite.Name)
	}
	if invite.Accepted {
		return nil
	}
	invite.Apply(Accepted{invite.Id})
	return nil
}

func (invite *Aggregate) Decline() error {
	if invite.Accepted {
		return fmt.Errorf("%s already accepted", invite.Name)
	}
	if invite.Declined {
		return nil
	}
	invite.Apply(Declined{invite.Id})
	return nil
}

func (invite *Aggregate) Apply(ev interface{}) {
	invite.Record(ev)
	switch ev := ev.(type) {
	case Accepted:
		invite.Accepted = true
	case Declined:
		invite.Declined = true
	case Created:
		invite.Name = ev.Name
	}
}
