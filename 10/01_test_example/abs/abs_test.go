package abs

import (
	"testing"
)

func Test_Abs_Positive(t *testing.T) {
	got := Abs(1)
	if got != 1 {
		t.Errorf("Abs(1) = %d; want 1", got)
	}
}

func Test_Abs_Negative(t *testing.T) {
	got := Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}
