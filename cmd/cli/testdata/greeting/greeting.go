package greeting

type Greeting int

const (
	GreetingWorld Greeting = iota
	GreetingMars
)

type InvalidNonContinuousGreeting int

const (
	InvalidNonContinuousGreetingWorld InvalidNonContinuousGreeting = iota
	_
	InvalidNonContinuousGreetingMars
)
