package handler

import "testing"

func TestFormatList(t *testing.T) {
	var original, expected, actual string

	original = "hello, world!"
	expected = "hello, world!"
	actual = formatList(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "* one"
	expected = "- one"
	actual = formatList(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "* one\n\n* two"
	expected = "- one\n\n- two"
	actual = formatList(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "  \n\n     * one    \n\n  * two"
	expected = "  \n\n     - one    \n\n  - two"
	actual = formatList(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "2*2"
	expected = "2*2"
	actual = formatList(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}
}

func TestFormatBold(t *testing.T) {
	var original, expected, actual string

	original = "hello, world!"
	expected = "hello, world!"
	actual = formatBold(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = `\*\*hello\*\*`
	expected = `**hello**`
	actual = formatBold(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}
}

func TestEscapeSpecialCharacters(t *testing.T) {
	var original, expected, actual string

	original = "hello, world!"
	expected = "hello, world!"
	actual = escapeSpecialCharacters(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "* `*`"
	expected = "\\* `*`"
	actual = escapeSpecialCharacters(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "```\n*\n```\n * `*` * `  *  ` ```*```*"
	expected = "```\n*\n```\n \\* `*` \\* `  *  ` ```*```\\*"
	actual = escapeSpecialCharacters(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}

	original = "```\nhello\n```\n\nworld"
	expected = "```\nhello\n```\nworld"
	actual = escapeSpecialCharacters(original)
	if expected != actual {
		t.Errorf("expected '%v', but received '%v'", expected, actual)
	}
}
