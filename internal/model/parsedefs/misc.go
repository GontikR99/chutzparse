package parsedefs

import "strconv"

// RenderAmount converts an amount into a brief string representation with at most 3 significant digits
func RenderAmount(amount float64) string {
	if amount < 0.9995 {
		return "   0 "
	} else if amount < 9.995 {
		intRep := strconv.FormatInt(int64(100*amount+0.5),10)
		return intRep[:1]+"."+intRep[1:]+" "
	} else if amount < 99.95 {
		intRep := strconv.FormatInt(int64(10*amount+0.5),10)
		return intRep[:2]+"."+intRep[2:]+" "
	} else if amount < 999.5 {
		intRep := strconv.FormatInt(int64(amount+0.5),10)
		return " "+intRep+" "
	} else if amount < 9999500 {
		return RenderAmount(amount/1000.0)[:4]+"k"
	} else {
		return RenderAmount(amount)[:4]+"M"
	}
}

const ColorLimeGreen = "#32CD32"
const ColorPastelRed = "#FAA0A0"
