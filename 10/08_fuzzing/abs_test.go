package abs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Fuzz_Abs(f *testing.F) {
	f.Add(42)
	f.Fuzz(func(t *testing.T, x int) {
		got := Abs(x)

		// prepare expected value
		expected := x
		if x < 0 {
			expected = -x
		}

		// uncomment this to introduce an error, then run regular go test -v ./...
		// if x > 100 {
		// 	expected = 213 // some non-sense number
		// }

		assert.Equal(t, expected, got)
	})
}

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
