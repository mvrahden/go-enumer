package project

// Note: This file serves to ensure, that types from various files are identified.

// AccountState represents various entity lifecycle states.
//go:enum -serializers=json,sql -transform=upper
type AccountState uint

const (
	AccountStateStaged AccountState = iota
	AccountStateProvisioned
	AccountStateActivated
	AccountStateDeactivated
	AccountStateDeprovisioned
)

// NotAnEnum does not contain the magic comment and will therefore be ignored.
type NotAnEnum uint
