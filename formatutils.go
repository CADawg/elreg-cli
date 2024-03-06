package main

import "strings"

func WrapText(text string) string {
	const maxLineLength = 80
	var result strings.Builder
	words := strings.Fields(text)
	currentLineLength := 0

	for _, word := range words {
		wordLength := len(word)
		if currentLineLength+wordLength+1 > maxLineLength && wordLength < maxLineLength {
			result.WriteString("\n")
			currentLineLength = 0
		}
		result.WriteString(word + " ")
		currentLineLength += wordLength + 1
	}

	return strings.TrimSpace(result.String())
}

func GetReadText(read bool) string {
	if read {
		return "[Read]"
	} else {
		return ""
	}
}
