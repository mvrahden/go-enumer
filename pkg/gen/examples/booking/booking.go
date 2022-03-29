package booking

// BookingState is an indicator for bookings.
//go:enumer -from=booking.csv -serializers=json -support=undefined
type BookingState int

const (
	BookingStateCreated BookingState = iota
	BookingStateUnavailable
	BookingStateFailed
	BookingStateCanceled
	BookingStateNotFound
	BookingStateDeleted
)
