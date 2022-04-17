package greetings

//go:enum
type Greeting uint8

const (
	GreetingРоссия Greeting = iota + 1
	Greeting中國
	Greeting日本
	Greeting한국
	GreetingČeskáRepublika
	Greeting𝜋
)

//go:enum
type GreetingWithDefault uint8

const (
	GreetingWithDefaultWorld GreetingWithDefault = iota
	GreetingWithDefaultРоссия
	GreetingWithDefault中國
	GreetingWithDefault日本
	GreetingWithDefault한국
	GreetingWithDefaultČeskáRepublika
	GreetingWithDefault𝜋
)
