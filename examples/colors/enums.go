package colors

import "fmt"

//go:enum -from=colors.csv
type Color uint

func (c Color) ToRGB() string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.GetRed(), c.GetGreen(), c.GetBlue())
}

func (c Color) ToRGBA() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%.3g)", c.GetRed(), c.GetGreen(), c.GetBlue(), c.GetAlpha())
}

func (c Color) ToHex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.GetRed(), c.GetGreen(), c.GetBlue())
}
