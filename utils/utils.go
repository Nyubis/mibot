package utils

import (
	"fmt"
	"regexp"
)

var whitespaceRe = regexp.MustCompile("\\s+")

func Contains(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}

func ParseSize(bytes int64) string {
	units := [...]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	size := float64(bytes)
	var i int

	for i = 0; size > 1024 && i < len(units)-1; i++ {
		size /= 1024
	}

	if i == 0 {
		// Decimal bytes aren't a thing, so the trailing zeroes we'd otherwise have look silly
		return fmt.Sprintf("%.0f B", size)
	} else {
		return fmt.Sprintf("%.2f %s", size, units[i])
	}
}

func Truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "..."
}

func NormalizeWhitespace(s string) string {
	return whitespaceRe.ReplaceAllString(s, " ")
}
