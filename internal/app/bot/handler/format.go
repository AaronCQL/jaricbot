package handler

import "strings"

// Formats the given text to be compatible with Telegram's 'Markdown' syntax.
// See notes in: https://core.telegram.org/bots/api#markdown-style.
func formatTelegramMarkdown(text string) string {
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

		// Escape special characters outside of code blocks
		if char == '`' {
			shouldEscape = !shouldEscape
		}
		if shouldEscape &&
			(char == '*' || char == '_' || char == '[' || char == ']') {
			sb.WriteRune('\\')
		}
		sb.WriteRune(char)
	}
	return sb.String()
}
