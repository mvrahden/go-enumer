package gen

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TypeParsing(t *testing.T) {
	type args struct {
		typ GoType
		raw string
	}
	tests := []struct {
		name string
		args args
		want any
		err  error
	}{
		{"succeeds", args{GoTypeComplex64, "7+2i"}, complex64(complex(7, 2)), nil},
		{"fails", args{GoTypeComplex64, "<!>"}, complex64(0), errors.New("some error")},
		{"succeeds", args{GoTypeComplex128, "7+2i"}, complex(7, 2), nil},
		{"fails", args{GoTypeComplex128, "<!>"}, complex128(0), errors.New("some error")},
		{"succeeds", args{GoTypeFloat32, "7.24568"}, float32(7.24568), nil},
		{"fails", args{GoTypeFloat32, "<!>"}, float32(0), errors.New("some error")},
		{"succeeds", args{GoTypeFloat64, "7.24568"}, float64(7.24568), nil},
		{"fails", args{GoTypeFloat64, "<!>"}, float64(0), errors.New("some error")},
		{"succeeds", args{GoTypeBool, "true"}, true, nil},
		{"succeeds", args{GoTypeBool, "1"}, true, nil},
		{"succeeds", args{GoTypeBool, "t"}, true, nil},
		{"succeeds", args{GoTypeBool, "false"}, false, nil},
		{"succeeds", args{GoTypeBool, "0"}, false, nil},
		{"succeeds", args{GoTypeBool, "f"}, false, nil},
		{"fails", args{GoTypeBool, "<!>"}, false, errors.New("some error")},
		{"succeeds", args{GoTypeUnsignedInteger, "123"}, uint(123), nil},
		{"fails", args{GoTypeUnsignedInteger, "<!>"}, uint(0), errors.New("some error")},
		{"succeeds", args{GoTypeUnsignedInteger8, "123"}, uint8(123), nil},
		{"fails", args{GoTypeUnsignedInteger8, "<!>"}, uint8(0), errors.New("some error")},
		{"succeeds", args{GoTypeUnsignedInteger16, "123"}, uint16(123), nil},
		{"fails", args{GoTypeUnsignedInteger16, "<!>"}, uint16(0), errors.New("some error")},
		{"succeeds", args{GoTypeUnsignedInteger32, "123"}, uint32(123), nil},
		{"fails", args{GoTypeUnsignedInteger32, "<!>"}, uint32(0), errors.New("some error")},
		{"succeeds", args{GoTypeUnsignedInteger64, "123"}, uint64(123), nil},
		{"fails", args{GoTypeUnsignedInteger64, "<!>"}, uint64(0), errors.New("some error")},
		{"succeeds", args{GoTypeSignedInteger, "123"}, int(123), nil},
		{"fails", args{GoTypeSignedInteger, "<!>"}, int(0), errors.New("some error")},
		{"succeeds", args{GoTypeSignedInteger8, "123"}, int8(123), nil},
		{"fails", args{GoTypeSignedInteger8, "<!>"}, int8(0), errors.New("some error")},
		{"succeeds", args{GoTypeSignedInteger16, "123"}, int16(123), nil},
		{"fails", args{GoTypeSignedInteger16, "<!>"}, int16(0), errors.New("some error")},
		{"succeeds", args{GoTypeSignedInteger32, "123"}, int32(123), nil},
		{"fails", args{GoTypeSignedInteger32, "<!>"}, int32(0), errors.New("some error")},
		{"succeeds", args{GoTypeSignedInteger64, "123"}, int64(123), nil},
		{"fails", args{GoTypeSignedInteger64, "<!>"}, int64(0), errors.New("some error")},
		{"succeeds", args{GoTypeString, "123"}, "123", nil},
		{"succeeds", args{GoTypeUnknown, "123"}, "123", nil},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Parse %T %s", tt.want, tt.name), func(t *testing.T) {
			fn := getParserFuncFor(tt.args.typ)
			v, err := fn(tt.args.raw)
			if tt.err != nil {
				require.NotZero(t, err)
				require.Zero(t, v)
				return
			}
			require.Zero(t, err)
			require.Equal(t, tt.want, v)
		})
	}
}
