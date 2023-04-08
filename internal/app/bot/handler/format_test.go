package handler

import "testing"

func TestFormatTelegramMarkdown(t *testing.T) {
	var original, expected, actual string

	original = "hello, world!"
	expected = "hello, world!"
	actual = formatTelegramMarkdown(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "* `*`"
	expected = "\\* `*`"
	actual = formatTelegramMarkdown(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "```\n*\n```\n * `*` * `  *  ` ```*```*"
	expected = "```\n*\n```\n \\* `*` \\* `  *  ` ```*```\\*"
	actual = formatTelegramMarkdown(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "```\nhello\n```\n\nworld"
	expected = "```\nhello\n```\nworld"
	actual = formatTelegramMarkdown(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}
}
