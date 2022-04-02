package project

// Note: This file serves to ensure, that types from various files are identified.

// State represents various entity lifecycle states.
//go:enumer -serializers=json,sql -transform=upper
type State uint

const (
	StateStaged State = iota
	StateProvisioned
	StateActivated
	StateDeactivated
	StateDeprovisioned
)

// NotAnEnum does not contain the magic comment and will therefore be ignored.
type NotAnEnum int
