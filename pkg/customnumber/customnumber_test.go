package customnumber_test

import (
	"github.com/slysterous/print-scrape/pkg/customnumber"
	"testing"
)

func TestCustomNumberString(t *testing.T) {
	want := "150000"
	values := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	number := customnumber.NewNumber(values, want)

	if got, want := number.String(), want; got != want {
		t.Errorf("String of custom number, want: %s got: %s", want, got)
	}
}

func TestIncrement(t *testing.T) {
	want := "150001"
	values := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	number := customnumber.NewNumber(values, "150000")
	number.Increment()
	if got, want := number.String(), want; got != want {
		t.Errorf("String of custom number, want: %s got: %s", want, got)
	}

}
