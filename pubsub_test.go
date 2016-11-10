package pubsub

import (
	"fmt"
	"sync"
	"testing"

	"github.com/paulormart/assert"
)

func TestHello(t *testing.T) {
	m := hello()
	assert.Equal(t, "hello", m)
}

type eventHello struct{}

func (e eventHello) String() string {
	return "event_hello"
}

func (e eventHello) Event() {
	fmt.Print("do something")
}

func TestSubscribe(t *testing.T) {
	wg := &sync.WaitGroup{}
	c := make(chan Event, 1)
	evh := &eventHello{}
	Subscribe(c, evh)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ev := <-c
		assert.Equal(t, "event_hello", ev.String())
	}()
	Publish(evh)
	wg.Wait()
}
