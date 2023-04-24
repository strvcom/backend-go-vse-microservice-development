package max

import (
	"testing"
)

func Test_Max(t *testing.T) {
	cases := []struct {
		args []int
		want int
	}{
		{
			args: []int{1, 2, 3},
			want: 3,
		},
		{
			args: []int{3, 2, 1},
			want: 3,
		},
		{
			args: []int{1, 2, 5},
			want: 5,
		},
	}
	for _, c := range cases {
		got := Max(c.args)
		if got != c.want {
			t.Errorf("Max(%v) = %d; want %d", c.args, got, c.want)
		}
	}
}
