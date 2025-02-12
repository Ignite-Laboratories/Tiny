package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_To_Number_TooWide(t *testing.T) {
	number := tiny.To.Number(222, 1, 0, 1, 0)
	if number != 10 {
		t.Errorf("Expected %d, Got %d", 10, number)
	}
}

func Test_To_Number_SameWidth(t *testing.T) {
	number := tiny.To.Number(4, 1, 0, 1, 0)
	if number != 10 {
		t.Errorf("Expected %d, Got %d", 5, number)
	}
}

func Test_To_Number_UnderWide(t *testing.T) {
	number := tiny.To.Number(3, 1, 0, 1, 0)
	if number != 5 {
		t.Errorf("Expected %d, Got %d", 5, number)
	}
}

func Test_To_Number_LargeNumber(t *testing.T) {
	number := tiny.To.Number(32, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1)
	if number != 1521080171 {
		t.Errorf("Expected %d, Got %d", 5, 1521080171)
	}
}
