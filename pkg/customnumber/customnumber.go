package customnumber

import (
	"bytes"
	"container/list"
	"container/ring"
)

// => 0 1 2 3 4 5 6 7 8 9 a b c d e f g h i j k l m n o p q r s t u v w x y z  arithmetic system

// var value = "0a9esd"

// value.increment() => 0a8ese

//==================================

// var value = "0a9esz"

// value.increment() => 0a8et0
//=====================================

//start scrapping
// 000000
// .
// .
// 00000z
// 000010

// Number represents a custom number.
type Number struct {
	Digits *list.List
}

// NewNumber initializes a CustomNumber list of x digits.
func NewNumber(values []rune, initial string) Number {
	// initialise a new number.
	number := Number{Digits: list.New()}
	// add digits to the number along with their state.
	for i := 0; i < len(initial); i++ {
		digit := newDigit(values, rune(initial[i]))
		number.Digits.PushBack(digit)
	}

	return number
}

// newDigit creates and initializes a new digit (ring).
func newDigit(values []rune, state rune) ring.Ring {
	// initialize a new empty ring
	r := ring.New(len(values))

	// fill the ring with values
	for _, e := range values {
		r.Value = e
		r = r.Next()
	}

	// roll the ring in desired "state" position.
	for range values {
		if r.Value == state {
			break
		}
		r = r.Next()
	}

	return *r
}

// Increment performs a +1 to the Number.
func (p *Number) Increment() {
	// take the second digit and keep going if there are any arithmetic holdings
	for e := p.Digits.Back(); e != nil; e = e.Prev() {
		r, ok := e.Value.(ring.Ring)
		if ok {
			// increment digit by one
			r = *r.Next()
			// update list item
			e.Value = r

			// if the digit is being reset then we
			// have an arithmetic holding
			if r.Value != '0' {
				return
			}
		}
	}
}

// String prints a string representation of Number.
func (p Number) String() string {
	// Loop over container list.
	var numberBytes bytes.Buffer
	for e := p.Digits.Front(); e != nil; e = e.Next() {
		r, ok := e.Value.(ring.Ring)
		if !ok {
			return ""
		}

		v, ok := r.Value.(rune)
		if !ok {
			return ""
		}

		numberBytes.WriteString(string(v))

	}
	return numberBytes.String()
}
