package membus

import (
	"math/rand"
	"testing"
	"time"

	"github.com/egonelbre/guestlist/event"
)

func randomEvent() Event {
	return int64(rand.Int())
}

func testCancellation(t *testing.T, bus event.Bus) {
	event := randomEvent()
	ok := true
	cancel := bus.Listen(func(e Event) { ok = false })
	cancel()
	bus.Publish(event)
	time.Sleep(50 * time.Microsecond)
	if !ok {
		t.Fatal("cancellation didn't work")
	}
}

func testPublish(t *testing.T, bus event.Bus) {
	event := randomEvent()
	done := make(chan int, 2)

	check := func(e Event) {
		if e != event {
			t.Fatalf("published event was different")
		}
		done <- 1
	}

	bus.Listen(check)
	bus.Listen(check)
	bus.Publish(event)
	<-done
	<-done
}

func TestBusCancellation(t *testing.T) {
	testCancellation(t, New())
}

func TestBusPublish(t *testing.T) {
	testPublish(t, New())
}
