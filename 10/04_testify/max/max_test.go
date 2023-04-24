package max

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Max(t *testing.T) {
	cases := []struct {
		name string
		args []int
		want int
	}{
		{
			name: "basic slice",
			args: []int{1, 2, 3},
			want: 3,
		},
		{
			name: "reverse slice",
			args: []int{3, 2, 1},
			want: 3,
		},
		{
			name: "random slice",
			args: []int{1, 2, 5},
			want: 5,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Max(c.args)
			assert.Equal(t, c.want, got)
		})
	}
}
