package max

import (
	"testing"
)

func Test_Max(t *testing.T) {
	got := Max([]int{1, 2, 3})
	if got != 3 {
		t.Errorf("Max([]int{1, 2, 3}) = %d; want 3", got)
	}
}
