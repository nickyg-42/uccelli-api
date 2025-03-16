package test

import (
	"nest/utils"
	"testing"
)

func TestIsSelfOrSA(t *testing.T) {
	userId := 1
	want := true
	result, err := utils.IsSelfOrSA(r, userId)
	if !want.MatchString(msg) || err != nil {
		t.Errorf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}
