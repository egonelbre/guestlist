package event

type Aggregate struct {
	Id      AggregateId
	Version int64
	Changes []Event
}

func (rec *Aggregate) Record(event Event)  { rec.Changes = append(rec.Changes, event) }
func (rec *Aggregate) GetId() AggregateId  { return rec.Id }
func (rec *Aggregate) GetVersion() int64   { return rec.Version }
func (rec *Aggregate) GetChanges() []Event { return rec.Changes }
