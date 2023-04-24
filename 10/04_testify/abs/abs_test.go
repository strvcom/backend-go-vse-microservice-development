package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, c.want, got)
		})
	}
}
