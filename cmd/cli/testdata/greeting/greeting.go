package greeting

//go:enumer
type Greeting int

const (
	GreetingWorld Greeting = iota
	GreetingMars
)
