package parsedefs

import "strconv"

// RenderAmount converts an amount into a brief string representation with at most 3 significant digits
func RenderAmount(amount float64) string {
	if amount < 999.5 {
		return strconv.FormatFloat(amount, 'g', 3, 64)
	} else if amount < 9999500 {
		amtFlt := float64(amount) / 1000.0
		return strconv.FormatFloat(amtFlt, 'g', 3, 64) + "k"
	} else {
		amtFlt := float64(amount) / 1000000.0
		return strconv.FormatFloat(amtFlt, 'g', 3, 64) + "M"
	}
}

// RenderFixed converts an amount into a brief string representation with exactly 3 significant digits, and exactly
// 5 characters.
func RenderFixed(amount float64) string {
	if amount < 0.09995 {
		return "  0. "
	} else if amount < 0.9995 {
		intRep := strconv.FormatInt(int64(1000*amount+0.5), 10)
		return "." + intRep + " "
	} else if amount < 9.995 {
		intRep := strconv.FormatInt(int64(100*amount+0.5), 10)
		return intRep[:1] + "." + intRep[1:] + " "
	} else if amount < 99.95 {
		intRep := strconv.FormatInt(int64(10*amount+0.5), 10)
		return intRep[:2] + "." + intRep[2:] + " "
	} else if amount < 999.5 {
		intRep := strconv.FormatInt(int64(amount+0.5), 10)
		return intRep + ". "
	} else if amount < 9999500 {
		return RenderFixed(amount / 1000.0)[:4] + "k"
	} else {
		return RenderFixed(amount)[:4] + "M"
	}
}

func RenderPercent(ratio float64) string {
	return RenderFixed(100 * ratio)[:4] + "%"
}

const ColorLimeGreen = "#32CD32"
const ColorPastelRed = "#FAA0A0"
