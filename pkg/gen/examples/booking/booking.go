package booking

// BookingState is an indicator for bookings.
//go:enumer -from=booking.csv -serializers=json -supports=undefined
type BookingState int

/* CSV-Source: BookingState
0,Created,The booking was created successfully
1,Unavailable,The booking was not available
2,Failed,The booking failed
3,Canceled,The booking was canceled
4,NotFound,The booking was not found
5,Deleted,The booking was deleted
*/

const (
	BookingStateCreated BookingState = iota
	BookingStateUnavailable
	BookingStateFailed
	BookingStateCanceled
	BookingStateNotFound
	BookingStateDeleted
)

////go:enumer
type Greeting uint8

const (
	GreetingРоссия Greeting = iota + 1
	Greeting中國
	Greeting日本
	Greeting한국
	GreetingČeskáRepublika
	Greeting𝜋
)

type PillSigned8 int8

const (
	PillSigned8Placebo PillSigned8 = iota
	PillSigned8Aspirin
	PillSigned8Ibuprofen
	PillSigned8Paracetamol
	PillSigned8Acetaminophen PillSigned8 = iota - 1
	PillSigned8VitaminC
)
