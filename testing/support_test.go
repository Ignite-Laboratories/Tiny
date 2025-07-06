package testing

import (
	"github.com/ignite-laboratories/tiny"
	"regexp"
	"testing"
	"unicode/utf8"
)

func Test_Support_PrintIndexWidth(t *testing.T) {
	// Separate regex for each output
	validateWithDigits := regexp.MustCompile(`^\|←\s*\d+\s*→\|$`)
	validateWithoutDigits := regexp.MustCompile(`^\|←\s*→\|$`)

	for i := 0; i < 129; i++ {
		// Test with width (with digits)
		strWith := tiny.PrintIndexWidth(i)
		strWithout := tiny.PrintIndexWidth(i, false)

		switch {
		case i == 0:
			if strWith != "||" && strWithout != "||" {
				t.Fatalf("expected '||', got %s", strWith)
			}
			continue
		case i == 1:
			if strWith != "|1|" && strWithout != "|1|" {
				t.Fatalf("expected '|1|', got %s", strWith)
			}
			continue
		case i == 2:
			if strWith != "|←→|" && strWithout != "|←→|" {
				t.Fatalf("expected '|←→|', got %s", strWith)
			}
			continue
		}

		length := utf8.RuneCountInString(strWith)
		if length != i+2 {
			t.Fatalf("strWith: expected %d characters between pipes, got %d", i+2, length)
		}
		length = utf8.RuneCountInString(strWithout)
		if length != i+2 {
			t.Fatalf("strWithout: expected %d characters between pipes, got %d", i+2, length)
		}

		if !validateWithDigits.MatchString(strWith) {
			t.Fatalf("strWith does not abstractly match the desired output '|← 𝑛 →|': %s", strWith)
		}
		if !validateWithoutDigits.MatchString(strWithout) {
			t.Fatalf("strWithout does not abstractly match the desired output '|←   →|': %s", strWith)
		}
	}
}

func Test_Support_PrintIndexWidth_NegativeInput(t *testing.T) {
	strWith := tiny.PrintIndexWidth(-1)
	strWithout := tiny.PrintIndexWidth(-1, false)

	if strWith != "||" && strWithout != "||" {
		t.Fatalf("expected '||', got %s", strWith)
	}
}
