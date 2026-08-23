package report

import "strconv"

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func Percent(part int, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(part) * 100 / float64(total)
}

func Plural(count int, singular string) string {
	if count == 1 {
		return singular
	}
	return singular + "s"
}
