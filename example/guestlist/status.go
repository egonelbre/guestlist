package main

import (
	"github.com/egonelbre/event"
	"github.com/egonelbre/event/example/guestlist/invitation"
)

type Status struct {
	Name  string
	State string
}

type StatusView struct {
	Invitations map[event.AggregateId]*Status
}

func NewStatusView() *StatusView {
	return &StatusView{make(map[event.AggregateId]*Status)}
}

func (sv *StatusView) Apply(ev event.Event) {
	switch ev := ev.(type) {
	case invitation.Created:
		sv.Invitations[ev.Id] = &Status{ev.Name, "unknown"}
	case invitation.Accepted:
		s := sv.Invitations[ev.Id]
		s.State = "accepted"
	case invitation.Declined:
		s := sv.Invitations[ev.Id]
		s.State = "declined"
	}
}
