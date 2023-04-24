package abs

import (
	"testing"
)

func Test_Abs(t *testing.T) {
	cases := []struct {
		name string
		args int
		want int
	}{
		{
			name: "positive",
			args: 1,
			want: 1,
		},
		{
			name: "negative",
			args: -1,
			want: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Abs(c.args)
			if got != c.want {
				t.Errorf("Abs(%d) = %d; want %d", c.args, got, c.want)
			}
		})
	}
}
