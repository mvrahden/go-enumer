package colors

import "fmt"

//go:enum -from=colors.csv
type Color uint

const (
	// ColorCyan represents Cyan.
	//go:enum assert={"6":"Cyan","7":"Magenta"}
	ColorCyan    Color = 6
	ColorMagenta Color = 7

	//go:enum assert={"0":"Black"}
	ColorBlack Color = 0
	//go:enum assert={"1":"White"}
	ColorWhite Color = 1
)

func (c Color) ToRGB() string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.GetRed(), c.GetGreen(), c.GetBlue())
}

func (c Color) ToRGBA() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%.3g)", c.GetRed(), c.GetGreen(), c.GetBlue(), c.GetAlpha())
}

func (c Color) ToHex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.GetRed(), c.GetGreen(), c.GetBlue())
}
