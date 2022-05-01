package colors

import (
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestEnums(t *testing.T) {
	t.Run("Color", func(t *testing.T) {
		type outcome struct {
			rgb  string
			rgba string
			hex  string
		}
		tc := []struct {
			raw    string
			wanted outcome
		}{
			{"black", outcome{"rgb(0,0,0)", "rgba(0,0,0,1)", "#000000"}},
			{"Black", outcome{"rgb(0,0,0)", "rgba(0,0,0,1)", "#000000"}},
			{"white", outcome{"rgb(255,255,255)", "rgba(255,255,255,1)", "#FFFFFF"}},
			{"White", outcome{"rgb(255,255,255)", "rgba(255,255,255,1)", "#FFFFFF"}},
			{"red", outcome{"rgb(255,0,0)", "rgba(255,0,0,1)", "#FF0000"}},
			{"Red", outcome{"rgb(255,0,0)", "rgba(255,0,0,1)", "#FF0000"}},
			{"lime", outcome{"rgb(0,255,0)", "rgba(0,255,0,1)", "#00FF00"}},
			{"Lime", outcome{"rgb(0,255,0)", "rgba(0,255,0,1)", "#00FF00"}},
			{"blue", outcome{"rgb(0,0,255)", "rgba(0,0,255,1)", "#0000FF"}},
			{"Blue", outcome{"rgb(0,0,255)", "rgba(0,0,255,1)", "#0000FF"}},
			{"yellow", outcome{"rgb(255,255,0)", "rgba(255,255,0,1)", "#FFFF00"}},
			{"Yellow", outcome{"rgb(255,255,0)", "rgba(255,255,0,1)", "#FFFF00"}},
			{"cyan", outcome{"rgb(0,255,255)", "rgba(0,255,255,1)", "#00FFFF"}},
			{"Cyan", outcome{"rgb(0,255,255)", "rgba(0,255,255,1)", "#00FFFF"}},
			{"magenta", outcome{"rgb(255,0,255)", "rgba(255,0,255,1)", "#FF00FF"}},
			{"Magenta", outcome{"rgb(255,0,255)", "rgba(255,0,255,1)", "#FF00FF"}},
			{"silver", outcome{"rgb(192,192,192)", "rgba(192,192,192,1)", "#C0C0C0"}},
			{"Silver", outcome{"rgb(192,192,192)", "rgba(192,192,192,1)", "#C0C0C0"}},
			{"gray", outcome{"rgb(128,128,128)", "rgba(128,128,128,1)", "#808080"}},
			{"Gray", outcome{"rgb(128,128,128)", "rgba(128,128,128,1)", "#808080"}},
			{"grey", outcome{"rgb(128,128,128)", "rgba(128,128,128,1)", "#808080"}},
			{"Grey", outcome{"rgb(128,128,128)", "rgba(128,128,128,1)", "#808080"}},
			{"maroon", outcome{"rgb(128,0,0)", "rgba(128,0,0,1)", "#800000"}},
			{"Maroon", outcome{"rgb(128,0,0)", "rgba(128,0,0,1)", "#800000"}},
			{"olive", outcome{"rgb(128,128,0)", "rgba(128,128,0,1)", "#808000"}},
			{"Olive", outcome{"rgb(128,128,0)", "rgba(128,128,0,1)", "#808000"}},
			{"green", outcome{"rgb(0,128,0)", "rgba(0,128,0,1)", "#008000"}},
			{"Green", outcome{"rgb(0,128,0)", "rgba(0,128,0,1)", "#008000"}},
			{"purple", outcome{"rgb(128,0,128)", "rgba(128,0,128,1)", "#800080"}},
			{"Purple", outcome{"rgb(128,0,128)", "rgba(128,0,128,1)", "#800080"}},
			{"teal", outcome{"rgb(0,128,128)", "rgba(0,128,128,1)", "#008080"}},
			{"Teal", outcome{"rgb(0,128,128)", "rgba(0,128,128,1)", "#008080"}},
			{"navy", outcome{"rgb(0,0,128)", "rgba(0,0,128,1)", "#000080"}},
			{"Navy", outcome{"rgb(0,0,128)", "rgba(0,0,128,1)", "#000080"}},
		}
		for _, c := range tc {
			actual := utils.Must(ColorFromStringIgnoreCase(c.raw))

			require.Equal(t, c.wanted.rgb, actual.ToRGB())
			require.Equal(t, c.wanted.rgba, actual.ToRGBA())
			require.Equal(t, c.wanted.hex, actual.ToHex())
		}
	})
}
