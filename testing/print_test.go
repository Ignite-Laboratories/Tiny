package testing

import (
	"github.com/ignite-laboratories/tiny"
	"regexp"
	"testing"
	"unicode/utf8"
)

func Test_Print_DeltaCharacter(t *testing.T) {
	a := tiny.Print.DeltaCharacter(0, 0)
	b := tiny.Print.DeltaCharacter(0, 1)
	c := tiny.Print.DeltaCharacter(1, 0)

	if a != "|" {
		t.Fatalf("expected '|', got '%s'", a)
	}
	if b != "\\" {
		t.Fatalf("expected '\\', got '%s'", a)
	}
	if c != "/" {
		t.Fatalf("expected '/', got '%s'", a)
	}
}

func Test_Print_IndexWidth(t *testing.T) {
	// Separate regex for each output
	validateWithDigits := regexp.MustCompile(`^\|←\s*\d+\s*→\|$`)
	validateWithoutDigits := regexp.MustCompile(`^\|←\s*→\|$`)

	for i := 0; i < 129; i++ {
		// Test with width (with digits)
		strWith := tiny.Print.IndexWidth(i)
		strWithout := tiny.Print.IndexWidth(i, false)

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

func Test_Print_IndexWidth_NegativeInput(t *testing.T) {
	strWith := tiny.Print.IndexWidth(-1)
	strWithout := tiny.Print.IndexWidth(-1, false)

	if strWith != "||" && strWithout != "||" {
		t.Fatalf("expected '||', got %s", strWith)
	}
}
