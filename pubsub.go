package pubsub

import "sync"

type handler struct {
	topics map[string]bool
}

func (h *handler) valid(topic string) bool {
	return h.topics[topic]
}

func (h *handler) set(topic string) {
	h.topics[topic] = true
}
func (h *handler) clear(topic string) {
	h.topics[topic] = false
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

// Publish publishes an event on the registered subscriber channels.
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

// Unsubscribe remove events from the map.
func Unsubscribe(c chan<- Event) {
	handlers.Lock()
	defer handlers.Unlock()

	_, ok := handlers.m[c]
	if !ok {
		return
	}

	// remove := func(n string) {
	// 	if n == "" {
	// 		return
	// 	}
	// 	if h.valid(n) {
	// 		handlers.ref[n]--
	// 		if handlers.ref[n] == 0 {
	// 			delete(handlers.ref, n)
	// 		}
	// 		h.clear(n)
	// 	}
	// }
	// remove(h.topics)
	delete(handlers.m, c)
}
