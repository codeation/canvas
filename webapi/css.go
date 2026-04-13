package webapi

import (
	"fmt"
)

func styleName(code int) string {
	switch code {
	case 0:
		return "normal"
	case 1:
		return "oblique"
	case 2:
		return "italic"
	default:
		return "normal"
	}
}

func variantName(code int) string {
	switch code {
	case 0:
		return "normal"
	case 1:
		return "small-caps"
	default:
		return "normal"
	}
}

func stretchName(code int) string {
	switch code {
	case 0:
		return "ultra-condensed"
	case 1:
		return "extra-condensed"
	case 2:
		return "condensed"
	case 3:
		return "semi-condensed"
	case 4:
		return "normal"
	case 5:
		return "semi-expanded"
	case 6:
		return "expanded"
	case 7:
		return "extra-expanded"
	case 8:
		return "ultra-expanded"
	default:
		return "normal"
	}
}

func fontValue(height int, style, variant, weight, stretch int, family string) string {
	return fmt.Sprintf("%s %s %d %s %dpx %s",
		styleName(style), variantName(variant), weight, stretchName(stretch), height, family)
}

func color(r, g, b, a uint16) string {
	return fmt.Sprintf("rgba(%d,%d,%d,%3.2f)", r/0x100, g/0x100, b/0x100, float64(a)/float64(0xFFFF))
}

func px(value int) string {
	return fmt.Sprintf("%dpx", value)
}
