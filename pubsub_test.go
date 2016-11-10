package pubsub

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/paulormart/assert"
)

type eventHello struct{}

func (e eventHello) String() string {
	return "event_hello"
}

func (e eventHello) Event() {
	fmt.Print("do something")
}

type eventBye struct{}

func (e eventBye) String() string {
	return "event_bye"
}

func (e eventBye) Event() {
	fmt.Print("do something")
}

func TestMultiEventsDifferentChannels(t *testing.T) {
	wg := &sync.WaitGroup{}
	c1 := make(chan Event, 1)
	evh := &eventHello{}
	c2 := make(chan Event, 1)
	evb := &eventBye{}
	assert.Equal(t, 0, len(handlers.m))
	assert.Equal(t, 0, len(handlers.ref))
	Subscribe(c1, evh)
	assert.Equal(t, 1, len(handlers.m))
	assert.Equal(t, int64(1), handlers.ref[evh.String()])
	assert.Equal(t, int64(0), handlers.ref[evb.String()])
	Subscribe(c2, evb)
	assert.Equal(t, 2, len(handlers.m))
	assert.Equal(t, int64(1), handlers.ref[evh.String()])
	assert.Equal(t, int64(1), handlers.ref[evb.String()])
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case ev1 := <-c1:
			assert.Equal(t, "event_hello", ev1.String())
		case ev2 := <-c2:
			assert.Equal(t, "event_bye", ev2.String())
		default:
		}
	}()
	// Publish
	Publish(evb)
	Publish(evh)
	wg.Wait()
	// Unsubscribe
	Unsubscribe(c1)
	assert.Equal(t, int64(0), handlers.ref[evh.String()])
	assert.Equal(t, int64(1), handlers.ref[evb.String()])
	Unsubscribe(c2)
	assert.Equal(t, int64(0), handlers.ref[evb.String()])
}

func TestMultiEventsSingleChannel(t *testing.T) {
	wg := &sync.WaitGroup{}
	c1 := make(chan Event, 1)
	evh := &eventHello{}
	evb := &eventBye{}
	assert.Equal(t, 0, len(handlers.m))
	assert.Equal(t, 0, len(handlers.ref))
	Subscribe(c1, evh, evb)
	assert.Equal(t, 1, len(handlers.m))
	assert.Equal(t, int64(1), handlers.ref[evh.String()])
	assert.Equal(t, int64(1), handlers.ref[evb.String()])
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case ev1 := <-c1:
			log.Print("passei aqui")
			assert.Equal(t, "event_bye", ev1.String())
		default:
		}
	}()
	// Publish
	// Publish(evb)
	Publish(evh)
	wg.Wait()
}
