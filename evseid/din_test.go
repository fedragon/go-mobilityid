package evseid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	dinInput = TestInput{
		CountryCode:   "+49",
		OperatorCode:  "810",
		PowerOutletId: "000*438",
	}
	expectedDinId    = &DinEvseId{newEvseId(dinInput.CountryCode, dinInput.OperatorCode, dinInput.PowerOutletId)}
	assertValidDinId = NewAssertValidId(dinInput.CountryCode, dinInput.OperatorCode, dinInput.PowerOutletId)
)

func TestDinEvseId_String(t *testing.T) {
	t.Run("returns a valid DIN string with separators", func(t *testing.T) {
		assert.Equal(t, "+49*810*000*438", expectedDinId.String())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, Reader, error)
	}{
		{
			name:          "parses a valid DIN DinEvseId",
			input:         "+49*810*000*438",
			runAssertions: assertValidDinId,
		},
		{
			name:          "parses a valid DIN DinEvseId without + in country code",
			input:         "49*810*000*438",
			runAssertions: assertValidDinId,
		},
		{
			name:          "returns an error for an DIN DinEvseId string with an invalid country code",
			input:         "AA*810*000*438",
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error for an invalid DIN DinEvseId string",
			input:         "XYZ",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := ParseDin(test.input)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewDinEvseId(t *testing.T) {
	cases := []struct {
		name                              string
		countryCode, operatorId, outletId string
		checkDigit                        rune
		runAssertions                     func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an DIN DinEvseId from its parts",
			countryCode:   dinInput.CountryCode,
			operatorId:    dinInput.OperatorCode,
			outletId:      dinInput.PowerOutletId,
			runAssertions: assertValidDinId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "AA",
			operatorId:    dinInput.OperatorCode,
			outletId:      dinInput.PowerOutletId,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "+12345",
			operatorId:    dinInput.OperatorCode,
			outletId:      dinInput.PowerOutletId,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   dinInput.CountryCode,
			operatorId:    "TNMA",
			outletId:      dinInput.PowerOutletId,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if outletId is too long",
			countryCode:   dinInput.CountryCode,
			operatorId:    dinInput.OperatorCode,
			outletId:      "0123456789012345678901234567890123456789",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewDinEvseId(test.countryCode, test.operatorId, test.outletId)

			test.runAssertions(t, id, err)
		})
	}
}
