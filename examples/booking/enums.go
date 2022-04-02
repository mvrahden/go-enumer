package booking

// BookingState is an indicator for bookings.
//go:enumer -from=booking.csv
type BookingState uint

const ()

// BookingStateWithConfig will have it's individual configuration.
//go:enumer -from=booking.csv -serializers=json,yaml -support=undefined
type BookingStateWithConfig uint

const ()

// BookingStateWithConstants will have a subset (compared to CSV source)
// of explicitly defined constants.
//go:enumer -from=booking.csv
type BookingStateWithConstants uint

const (
	BookingStateWithConstantsCreated     BookingStateWithConstants = 0
	BookingStateWithConstantsUnavailable                           = 1
	BookingStateWithConstantsCanceled                              = 3
	BookingStateWithConstantsDeleted                               = 5
)
