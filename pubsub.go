package pubsub

import "sync"

type handler struct {
	events map[string]Event
}

func (h *handler) valid(e Event) bool {
	_, ok := h.events[e.String()]
	return ok
}

func (h *handler) set(e Event) {
	h.events[e.String()] = e
}

func (h *handler) clear(e Event) {
	delete(h.events, e.String())
}

// handlers is a collection of events.
var handlers struct {
	sync.Mutex
	m   map[chan<- Event]*handler
	ref map[string]int64
}

// Subscribe causes package pubsub to relay incoming events to c.
// If no events are provided, all incoming events will be relayed to c.
// Otherwise, just the provided events will.
//
// Package pubsub will not block sending to c: the caller must ensure
// that c has sufficient buffer space to keep up with the expected
// event rate. For a channel used for notification of just one event value,
// a buffer of size 1 is sufficient.
//
// It is allowed to call Subscribe multiple times with the same channel:
// each call expands the set of events sent to that channel.
//
// It is allowed to call Subscribe multiple times with different channels
// and the same events: each channel receives copies of incoming
// events independently.
func Subscribe(c chan<- Event, events ...Event) {
	if c == nil {
		panic("pubsub: subscribe using nil channel")
	}

	handlers.Lock()
	defer handlers.Unlock()

	h, ok := handlers.m[c]
	if !ok {
		if handlers.m == nil {
			handlers.m = make(map[chan<- Event]*handler)
		}
		if handlers.ref == nil {
			handlers.ref = make(map[string]int64)
		}
		h = &handler{make(map[string]Event)}
		handlers.m[c] = h
	}

	add := func(e Event) {
		if e == nil || e.String() == "" {
			return
		}
		if !h.valid(e) {
			h.set(e)
			handlers.ref[e.String()]++
		}
	}

	for _, e := range events {
		add(e)
	}
}

// Publish publishes an event on the registered subscriber channels.
func Publish(e Event) {
	if e == nil || e.String() == "" {
		return
	}

	handlers.Lock()
	defer handlers.Unlock()

	for c, h := range handlers.m {
		if h.valid(e) {
			// send but do not block for it
			select {
			case c <- e:
			default:
			}
		}
	}
}

// Unsubscribe remove events from the map.
func Unsubscribe(c chan<- Event) {
	handlers.Lock()
	defer handlers.Unlock()

	h, ok := handlers.m[c]
	if !ok {
		return
	}

	remove := func(e Event) {
		if e == nil || e.String() == "" {
			return
		}
		if h.valid(e) {
			handlers.ref[e.String()]--
			if handlers.ref[e.String()] == 0 {
				delete(handlers.ref, e.String())
			}
			h.clear(e)
		}
	}

	for _, e := range h.events {
		remove(e)
	}

	delete(handlers.m, c)
}
