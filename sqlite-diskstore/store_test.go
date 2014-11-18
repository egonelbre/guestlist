package diskstore

import (
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/egonelbre/event"
	"github.com/egonelbre/event/membus"
)

func TestStore(t *testing.T) {
	os.Remove("test.db")

	pub := membus.New()
	store, err := New("test.db", pub)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	a, a0, a1 := event.GenerateId(), 1, 2
	b, b0, b1 := event.GenerateId(), 3, 4

	checke := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}

	// save everything to disk

	checke(store.Save(a, -1, a0))
	checke(store.Save(b, -1, b0, b1))
	checke(store.Save(a, -1, a1))

	store.Close()

	// try to load the same file from disk
	store, err = New("test.db", pub)
	if err != nil {
		t.Fatalf("failed to open again db: %v", err)
	}
	defer store.Close()

	check := func(ok bool, i int) {
		if !ok {
			t.Fatalf("%d check failed", i)
		}
	}

	step := 0
	pub.Listen(func(e event.Event) {
		t.Logf("got %v", e)
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

	if err := store.Load(); err != nil {
		log.Fatalf("loading failed: %v", err)
	}

	aevents, ok := store.List(a)
	check(ok &&
		aevents[0].Id == a && aevents[0].Data == a0 &&
		aevents[1].Id == a && aevents[1].Data == a1, 5)
	bevents, ok := store.List(b)
	check(ok &&
		bevents[0].Id == b && bevents[0].Data == b0 &&
		bevents[1].Id == b && bevents[1].Data == b1, 6)
}

func TestStoreProject(t *testing.T) {
	os.Remove("test2.db")
	defer os.Remove("test2.db")

	pub := membus.New()
	store, err := New("test2.db", pub)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	var total int64
	uuid := event.GenerateId()
	N := int64(rand.Intn(100) + 50)
	for i := int64(0); i < N; i += 1 {
		v := int64(rand.Intn(256))
		total = (total ^ v) << 3
		if err != store.Save(uuid, i, v) {
			t.Fatal(err)
		}
	}
	store.Close()

	store, err = New("test2.db", pub)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	project := int64(0)
	pub.Listen(func(e event.Event) {
		v := e.(int64)
		project = (project ^ v) << 3
	})
	store.Load()
	store.Close()

	if total != project {
		t.Fatalf("something went wrong")
	}
}

func TestStoreConcurrencyError(t *testing.T) {
	pub := membus.New()
	store, err := New("test3.db", pub)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer os.Remove("test3.db")
	defer store.Close()

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
