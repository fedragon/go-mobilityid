package iso

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompute(t *testing.T) {
	cases := []struct {
		input              string
		expectedErr        bool
		expectedCheckDigit rune
	}{
		{
			input:              "NN123ABCDEFGHI",
			expectedErr:        false,
			expectedCheckDigit: 'T',
		},
		{
			input:              "FRXYZ123456789",
			expectedErr:        false,
			expectedCheckDigit: '2',
		},
		{
			input:              "ITA1B2C3E4F5G6",
			expectedErr:        false,
			expectedCheckDigit: '4',
		},
		{
			input:              "ESZU8WOX834H1D",
			expectedErr:        false,
			expectedCheckDigit: 'R',
		},
		{
			input:              "PT73902837ABCZ",
			expectedErr:        false,
			expectedCheckDigit: 'Z',
		},
		{
			input:              "DE83DUIEN83QGZ",
			expectedErr:        false,
			expectedCheckDigit: 'D',
		},
		{
			input:              "DE83DUIEN83ZGQ",
			expectedErr:        false,
			expectedCheckDigit: 'M',
		},
		{
			input:              "DE8AA001234567",
			expectedErr:        false,
			expectedCheckDigit: '0',
		},
		{
			input:              "Европарулит123",
			expectedErr:        true,
			expectedCheckDigit: -1,
		},
		{
			input:              "DE٨٣DUIEN٨٣QGZ",
			expectedErr:        true,
			expectedCheckDigit: -1,
		},
		{
			input:              "Å∏@*(Td\\uD83D\\uDE3BgaR^&(%",
			expectedErr:        true,
			expectedCheckDigit: -1,
		},
		{
			input:              "Å∏@*(Td\\uD83D\\uDE3BgR^&(%",
			expectedErr:        true,
			expectedCheckDigit: -1,
		},
		{
			input:              "",
			expectedErr:        true,
			expectedCheckDigit: -1,
		},
		{
			input:              "DE8AA0012345678",
			expectedErr:        true,
			expectedCheckDigit: -1,
		},
	}

	for _, testCase := range cases {
		result, err := ComputeCheckDigit(testCase.input)

		if err != nil && !testCase.expectedErr {
			t.Errorf("expected result but got error: %v", err)
		}

		assert.Equal(t, testCase.expectedCheckDigit, result)
	}
}
