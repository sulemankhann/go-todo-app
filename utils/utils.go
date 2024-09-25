package utils

import "strings"

// WrapText wraps the given text into lines of max width
func WrapText(text string, maxWidth int) (result string) {
	words := strings.Split(text, " ")
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > maxWidth {
			result += line + "\n"
			line = word
		} else {
			if len(line) > 0 {
				line += " "
			}
			line += word
		}
	}

	result += line

	return
}
