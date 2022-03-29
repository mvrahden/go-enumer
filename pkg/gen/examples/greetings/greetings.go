package greetings

//go:enumer
type Greeting uint8

const (
	GreetingРоссия Greeting = iota + 1
	Greeting中國
	Greeting日本
	Greeting한국
	GreetingČeskáRepublika
	Greeting𝜋
)
