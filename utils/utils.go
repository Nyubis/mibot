package utils

import (
	"fmt"
)

func Contains(hay []string, needle string) bool {
	for _, s := range hay {
		if s == needle {
			return true
		}
	}
	return false
}

func ParseSize(bytes int64) string {
	units := [...]string{ "B", "KiB", "MiB", "GiB", "TiB", "PiB" }
	size := float64(bytes)
	var i int

	for i = 0; size > 1024 && i < len(units)-1; i++ {
		size /= 1024
	}

	return fmt.Sprintf("%.2f %s", size, units[i])
}

func Truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "..."
}


