package booking

// BookingState is an indicator for bookings.
//go:enum -from=booking.csv
type BookingState uint

// BookingStateWithConfig will have its own configuration.
//go:enum -from=booking.csv -serializers=json,yaml -support=undefined
type BookingStateWithConfig uint

// BookingStateWithConstants will have a subset (compared to CSV source)
// of explicitly defined constants.
//go:enum -from=booking.csv
type BookingStateWithConstants uint

// With CSV sources, you can add your very own subset of constants.
// This can be handy if you often need to refer to very specific values
// of a large set of enums.
// Just add them here as constants.
const (
	BookingStateWithConstantsCreated     BookingStateWithConstants = 0
	BookingStateWithConstantsUnavailable                           = 1
	BookingStateWithConstantsCanceled                              = 3
	BookingStateWithConstantsDeleted                               = 5
)
