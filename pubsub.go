package pubsub

import (
	"fmt"
	"log"
	"sync"
)

func hello() string {
	return fmt.Sprintf("hello")
}

// A Event represents an message event.
type Event interface {
	String() string // event name
	Event()         // callback event method
}

type handler struct {
	name string
}

func (h *handler) valid(name string) bool {
	return h.name != ""
}

func (h *handler) set(name string) {
	h.name = name
}

// PubSub is a collection of events.
var handlers struct {
	sync.Mutex
	m   map[chan<- Event]*handler
	ref map[string]int64
}

// Subscribe ...
func Subscribe(c chan<- Event, events ...Event) {
	if c == nil {
		panic("pubsub: subscribe using nil channel")
	}

	handlers.Lock()
	defer handlers.Unlock()

	h, ok := handlers.m[c]
	log.Printf("handlers.cache[c] h:%v ok:%v", h, ok)
	if !ok {
		if handlers.m == nil {
			handlers.m = make(map[chan<- Event]*handler)
		}
		if handlers.ref == nil {
			handlers.ref = make(map[string]int64)
		}
		h = new(handler)
		handlers.m[c] = h
	}

	add := func(n string) {
		if n == "" {
			return
		}
		if !h.valid(n) {
			h.set(n)
			handlers.ref[n]++
		}
	}

	for _, ev := range events {
		add(ev.String())
	}
}

// Publish ...
func Publish(ev Event) {
	if ev.String() == "" {
		return
	}
	handlers.Lock()
	defer handlers.Unlock()

	for c, h := range handlers.m {
		if h.valid(ev.String()) {
			// send but do not block for it
			select {
			case c <- ev:
			default:
			}
		}
	}
}

// Unsubscribe ...
func Unsubscribe() {}
