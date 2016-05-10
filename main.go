package main

import (
	"log"

	"github.com/egonelbre/guestlist/event"
	"github.com/egonelbre/guestlist/event/membus"
	"github.com/egonelbre/guestlist/event/memstore"

	"github.com/egonelbre/guestlist/invitation"
)

var (
	bus   event.Bus
	store event.Store

	service *invitation.Service
	counter *CounterView
	status  *StatusView
)

func main() {
	bus = membus.New()
	store = memstore.New(bus)

	service = invitation.NewService(bus, store)

	counter = &CounterView{}
	bus.Listen(counter.Apply)

	status = NewStatusView()
	bus.Listen(status.Apply)

	log.Printf("Before: %#+v\n", counter)

	athenaId, err := service.NewInvite("Athena")
	check(err)

	check(service.AcceptInvite(athenaId))
	check(service.DeclineInvite(athenaId))

	hadesId, err := service.NewInvite("Hades")
	check(err)
	check(service.DeclineInvite(hadesId))

	_, err = service.NewInvite("Zeus")
	check(err)

	log.Printf("After: %#+v\n", counter)
	log.Printf("Status:\n")
	for _, s := range status.Invitations {
		log.Printf("  %s\t%s\n", s.Name, s.State)
	}
}

func check(err error) {
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}
