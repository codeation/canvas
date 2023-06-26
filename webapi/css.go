package webapi

import (
	"fmt"
)

var styleName = map[int]string{
	0: "normal",
	1: "oblique",
	2: "italic",
}

var variantName = map[int]string{
	0: "normal",
	1: "small-caps",
}

var stretchName = map[int]string{
	0: "ultra-condensed",
	1: "extra-condensed",
	2: "condensed",
	3: "semi-condensed",
	4: "normal",
	5: "semi-expanded",
	6: "expanded",
	7: "extra-expanded",
	8: "ultra-expanded",
}

func fontValue(height int, style, variant, weight, stretch int, family string) string {
	return fmt.Sprintf("%s %s %d %s %dpx %s",
		styleName[style], variantName[variant], weight, stretchName[stretch], height, family)
}

func color(r, g, b, a uint16) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%3.2f)", r/0x100, g/0x100, b/0x100, float64(a)/float64(0xFFFF))
}

func px(value int) string {
	return fmt.Sprintf("%dpx", value)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
