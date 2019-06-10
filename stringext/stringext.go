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
			splits = append(splits, sStr)
			if index < len(element) {
				eStr := element[index:]
				splits = append(splits, eStr)
			}
		} else {
			splits = append(splits, element)
		}
	}
	return splits
}
