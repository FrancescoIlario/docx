package stringext

import "strings"

// SplitAfterWithSeparator splits a string by the mean of separator and returns a slice containing the entries for separator
func SplitAfterWithSeparator(toSplit, sep string) []string {
	var splits []string
	s := strings.SplitAfter(toSplit, sep)
	for _, element := range s {
		index := strings.Index(element, sep)

		if index >= 0 {
			sStr := element[0:index]
			splits = appendIfNotEmpty(splits, sStr)

			if index < len(element) {
				eStr := element[index:]
				splits = appendIfNotEmpty(splits, eStr)
			}
		} else {
			splits = appendIfNotEmpty(splits, element)
		}
	}
	return splits
}

func appendIfNotEmpty(to []string, nw string) []string {
	if nw != "" {
		return append(to, nw)
	}
	return to
}
