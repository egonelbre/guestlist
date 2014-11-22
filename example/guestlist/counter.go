package main

import (
	"github.com/egonelbre/event"
	"github.com/egonelbre/event/example/guestlist/invitation"
)

type CounterView struct {
	NumUnknown  int
	NumAccepted int
	NumDeclined int
}

func (c *CounterView) Apply(ev event.Event) {
	switch ev.(type) {
	case invitation.Created:
		c.NumUnknown += 1
	case invitation.Accepted:
		c.NumUnknown -= 1
		c.NumAccepted += 1
	case invitation.Declined:
		c.NumUnknown -= 1
		c.NumDeclined += 1
	}
}
