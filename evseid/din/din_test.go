package din

import (
	"github.com/stretchr/testify/assert"
	"mobilityid/evseid"
	"testing"
)

var (
	input = evseid.TestInput{
		CountryCode:   "+49",
		OperatorCode:  "810",
		PowerOutletId: "000*438",
	}
	expectedId    = &EvseId{evseid.Id(input.CountryCode, input.OperatorCode, input.PowerOutletId)}
	assertValidId = evseid.NewAssertValidId(input.CountryCode, input.OperatorCode, input.PowerOutletId)
)

func TestEvseId_String(t *testing.T) {
	t.Run("returns a valid DIN string with separators", func(t *testing.T) {
		assert.Equal(t, "+49*810*000*438", expectedId.String())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, evseid.Reader, error)
	}{
		{
			name:          "parses a valid DIN EvseId",
			input:         "+49*810*000*438",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid DIN EvseId without + in country code",
			input:         "49*810*000*438",
			runAssertions: assertValidId,
		},
		{
			name:          "returns an error for an DIN EvseId string with an invalid country code",
			input:         "AA*810*000*438",
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error for an invalid DIN EvseId string",
			input:         "XYZ",
			runAssertions: evseid.AssertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := Parse(test.input)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewEvseId(t *testing.T) {
	cases := []struct {
		name                              string
		countryCode, operatorId, outletId string
		checkDigit                        rune
		runAssertions                     func(*testing.T, evseid.Reader, error)
	}{
		{
			name:          "creates an DIN EvseId from its parts",
			countryCode:   input.CountryCode,
			operatorId:    input.OperatorCode,
			outletId:      input.PowerOutletId,
			runAssertions: assertValidId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "AA",
			operatorId:    input.OperatorCode,
			outletId:      input.PowerOutletId,
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "+12345",
			operatorId:    input.OperatorCode,
			outletId:      input.PowerOutletId,
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   input.CountryCode,
			operatorId:    "TNMA",
			outletId:      input.PowerOutletId,
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error if outletId is too long",
			countryCode:   input.CountryCode,
			operatorId:    input.OperatorCode,
			outletId:      "0123456789012345678901234567890123456789",
			runAssertions: evseid.AssertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewEvseId(test.countryCode, test.operatorId, test.outletId)

			test.runAssertions(t, id, err)
		})
	}
}
