package validator

import "testing"

func TestRange_Not_IsValid(t *testing.T) {

	if !NewRange(10, 15).Value(12).IsValid() {
		t.Error("integer comparison error")
	}

	if !NewRange(0.9, 0.99).Value(0.95).IsValid() {
		t.Error("float comparison error")
	}

	if !NewRange("a", "c").Value("b").IsValid() {
		t.Error("string comparison error")
	}
}

func TestRange_IsValid(t *testing.T) {

	if NewRange(10, 15).Value(20).IsValid() {
		t.Error("integer comparison error")
	}

	if NewRange(0.9, 0.99).Value(1).IsValid() {
		t.Error("float comparison error")
	}

	if NewRange("a", "c").Value("z").IsValid() {
		t.Error("string comparison error")
	}
}
