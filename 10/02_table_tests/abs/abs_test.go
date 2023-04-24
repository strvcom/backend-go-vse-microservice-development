package abs

import (
	"testing"
)

func Test_Abs(t *testing.T) {
	cases := []struct {
		args int
		want int
	}{
		{
			args: 1,
			want: 1,
		},
		{
			args: -1,
			want: 1,
		},
	}

	for _, c := range cases {
		got := Abs(c.args)
		if got != c.want {
			t.Errorf("Abs(%d) = %d; want %d", c.args, got, c.want)
		}
	}
}
