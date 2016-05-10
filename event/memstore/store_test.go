package memstore

import (
	"math/rand"
	"reflect"
	"sync"
	"testing"

	"github.com/egonelbre/guestlist/event"
	"github.com/egonelbre/guestlist/event/membus"
)

func TestStore(t *testing.T) {
	pub := membus.New()
	store := New(pub)

	a, a0, a1 := event.GenerateId(), 1, 2
	b, b0, b1 := event.GenerateId(), 3, 4

	check := func(ok bool, i int) {
		if !ok {
			t.Fatalf("%d check failed", i)
		}
	}

	checke := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}

	step := 0
	pub.Listen(func(e event.Event) {
		switch step {
		case 0:
			check(reflect.DeepEqual(e, 1), 1)
		case 1:
			check(reflect.DeepEqual(e, 3), 2)
		case 2:
			check(reflect.DeepEqual(e, 4), 3)
		case 3:
			check(reflect.DeepEqual(e, 2), 4)
		}
		step += 1
	})

	checke(store.Save(a, -1, a0))
	checke(store.Save(b, -1, b0, b1))
	checke(store.Save(a, -1, a1))

	aevents, ok := store.List(a)
	check(ok &&
		aevents[0].Id == a && aevents[0].Data == a0 &&
		aevents[1].Id == a && aevents[1].Data == a1, 4)
	bevents, ok := store.List(b)
	check(ok &&
		bevents[0].Id == b && bevents[0].Data == b0 &&
		bevents[1].Id == b && bevents[1].Data == b1, 5)
}

func TestStoreProject(t *testing.T) {
	pub := membus.New()
	store := New(pub)

	uuid := event.GenerateId()

	m := sync.Mutex{}
	project := int64(0)

	pub.Listen(func(e event.Event) {
		v := e.(int64)
		m.Lock()
		project = (project ^ v) << 3
		m.Unlock()
	})

	total := int64(0)

	for i := 0; i < rand.Intn(100)+50; i += 1 {
		v := int64(rand.Intn(256))
		total = (total ^ v) << 3
		store.Save(uuid, int64(i), v)
	}

	if total != project {
		t.Fatalf("something went wrong")
	}
}

func TestStoreConcurrencyError(t *testing.T) {
	pub := membus.New()
	store := New(pub)
	uuid := event.GenerateId()

	if err := store.Save(uuid, 0, 1); err != nil {
		t.Fatalf("didn't expect error %v", err)
	}

	if store.Save(uuid, 0, 2) != event.ConcurrencyError {
		t.Fatalf("should have concurrency error")
	}

	if err := store.Save(uuid, -1, 3); err != nil {
		t.Fatalf("didn't expect error %v", err)
	}

	items, ok := store.List(uuid)
	if !ok || items[0].Data != 1 || items[1].Data != 3 {
		t.Fatalf("result event stream is wrong")
	}
}
