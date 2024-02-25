package handler

import (
	"regexp"
	"strings"
)

var (
	asteriskListRegex  = regexp.MustCompile(`(?m)^([ ]*)\*[ ]`)
	escapedBoldRegex   = regexp.MustCompile(`(\\\*\\\*)`)
	escapedItalicRegex = regexp.MustCompile(`(\\\_\\\_)`)
)

func formatList(text string) string {
	return asteriskListRegex.ReplaceAllString(text, "$1- ")
}

func escapeSpecialCharacters(text string) string {
	shouldEscape := true
	sb := strings.Builder{}
	for i, char := range text {
		// Ensure only one newline after triple backticks
		if i >= 2 &&
			text[i-2] == '`' &&
			text[i-1] == '\n' &&
			text[i] == '\n' {
			continue
		}

		if char == '`' {
			shouldEscape = !shouldEscape
		}
		if shouldEscape &&
			(char == '*' || char == '_') {
			sb.WriteRune('\\')
		}
		sb.WriteRune(char)
	}
	return sb.String()
}

func formatBold(text string) string {
	return escapedBoldRegex.ReplaceAllString(text, "*")
}

func formatItalic(text string) string {
	return escapedItalicRegex.ReplaceAllString(text, "_")
}

// Formats the given text to be compatible with Telegram's 'Markdown' syntax.
// See notes in: https://core.telegram.org/bots/api#markdown-style.
func formatTelegramMarkdown(text string) string {
	return strings.TrimSpace(
		formatItalic(
			formatBold(
				escapeSpecialCharacters(
					formatList(text),
				),
			),
		),
	)
}
