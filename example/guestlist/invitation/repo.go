package invitation

import "github.com/egonelbre/event"

type Repository struct {
	event.Store
}

func (repo *Repository) Create() *Aggregate {
	invite := &Aggregate{}
	invite.Id = event.GenerateId()
	return invite
}

func (repo *Repository) ById(id event.AggregateId) (invite *Aggregate, ok bool) {
	events, ok := repo.Store.List(id)
	if !ok {
		return nil, false
	}
	invite = &Aggregate{}
	invite.Id = id
	for _, event := range events {
		invite.Version = event.Version
		invite.Apply(event.Data)
	}
	invite.Changes = nil
	return invite, true
}
