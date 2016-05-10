package invitation

import (
	"fmt"

	"github.com/egonelbre/guestlist/event"
)

type Service struct {
	bus   event.Bus
	store event.Store

	repo *Repository
}

func NewService(bus event.Bus, store event.Store) *Service {
	return &Service{
		bus:   bus,
		store: store,

		repo: &Repository{store},
	}
}

func (srv *Service) NewInvite(name string) (event.AggregateId, error) {
	invite := srv.repo.Create()
	invite.Init(name)
	return invite.Id, srv.repo.SaveChanges(invite)
}

func (srv *Service) AcceptInvite(id event.AggregateId) error {
	invite, ok := srv.repo.ById(id)
	if !ok {
		return fmt.Errorf("invitation [%s] does not exist", id)
	}
	err := invite.Accept()
	if err != nil {
		return err
	}
	return srv.repo.SaveChanges(invite)
}

func (srv *Service) DeclineInvite(id event.AggregateId) error {
	invite, ok := srv.repo.ById(id)
	if !ok {
		return fmt.Errorf("invitation [%s] does not exist", id)
	}
	err := invite.Decline()
	if err != nil {
		return err
	}
	return srv.repo.SaveChanges(invite)
}
