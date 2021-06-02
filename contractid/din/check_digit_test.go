package din

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompute(t *testing.T) {
	cases := []struct {
		input    string
		expected rune
	}{
		{
			input:    "INTNM000071",
			expected: '9',
		},
		{
			input:    "INTNM000110",
			expected: 'X',
		},
		{
			input:    "INTNM000124",
			expected: '0',
		},
		{
			input:    "INTNM000114",
			expected: '6',
		},
		{
			input:    "INTNM000191",
			expected: '5',
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.expected, ComputeCheckDigit(testCase.input))
	}
}
