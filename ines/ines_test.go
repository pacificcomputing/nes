package ines

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("../tests/cpu/branch_timing_tests/1.Branch_Basics.nes")
	if err != nil {
		t.Errorf("Error opening test case: %v", err)
	}
	defer f.Close()
	if _, err := Parse(f); err != nil {
		t.Errorf("Error parsing header: %v", err)
	}
}
