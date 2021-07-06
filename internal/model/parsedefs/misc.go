package parsedefs

import "strconv"

// RenderAmount converts an amount into a brief string representation with at most 3 significant digits
func RenderAmount(amount float64) string {
	if amount < 1000 {
		return strconv.FormatFloat(amount, 'g', 3, 64)
	} else if amount < 1000000 {
		amtFlt := float64(amount) / 1000.0
		return strconv.FormatFloat(amtFlt, 'g', 3, 64) + "k"
	} else {
		amtFlt := float64(amount) / 1000000.0
		return strconv.FormatFloat(amtFlt, 'g', 3, 64) + "M"
	}
}

const ColorLimeGreen = "#32CD32"
const ColorPastelRed = "#FAA0A0"
