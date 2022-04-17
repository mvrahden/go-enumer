package greeting

//go:enum
type Greeting int

const (
	GreetingWorld Greeting = iota
	GreetingMars
)
