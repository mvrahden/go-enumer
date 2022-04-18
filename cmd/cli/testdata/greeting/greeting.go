package greeting

//go:enum
type Greeting uint

const (
	GreetingWorld Greeting = iota
	GreetingMars
)
