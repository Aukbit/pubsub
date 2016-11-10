package pubsub

// A Event represents an message event.
type Event interface {
	String() string // event name
	Event()         // callback event method
}
