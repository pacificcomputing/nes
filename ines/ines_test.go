package ines

import (
	"os"
	"testing"
)

func TestHeader_Read(t *testing.T) {
	f, err := os.Open("../tests/cpu/branch_timing_tests/1.Branch_Basics.nes")
	if err != nil {
		t.Errorf("Error opening test case: %v", err)
	}
	h := Header{}
	if err := h.Read(f); err != nil {
		t.Errorf("Error parsing header: %v", err)
	}
}
