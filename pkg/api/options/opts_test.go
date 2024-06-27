package options

import "testing"

type options int

const (
	_ options = iota
	OptionOne
	OptionTwo
)

func TestHas(t *testing.T) {

	all := []options{OptionOne, OptionTwo}
	one := []options{OptionOne}

	if !Has(OptionOne, all) {
		t.Errorf("OptionOne should exist in %+v", all)
	}

	if Has(OptionTwo, one) {
		t.Errorf("OptionTwo should not exist in %+v", one)
	}

	if Has(OptionOne, []options{}) {
		t.Errorf("OptionOne should not exist in %+v", []options{})
	}
}
