package iso

import (
	"github.com/stretchr/testify/assert"
	"mobilityid/evseid"
	"testing"
)

var (
	input = evseid.TestInput{
		CountryCode:   "DE",
		OperatorCode:  "AB7",
		PowerOutletId: "840*6487",
	}
	expectedId           = &EvseId{evseid.Id(input.CountryCode, input.OperatorCode, input.PowerOutletId)}
	assertValidId        = evseid.NewAssertValidId(input.CountryCode, input.OperatorCode, input.PowerOutletId)
	assertValidCompactId = evseid.NewAssertValidCompactId(input.CountryCode, input.OperatorCode, input.PowerOutletId)
)

func TestEvseId_String(t *testing.T) {
	t.Run("returns a valid ISO string with separators", func(t *testing.T) {
		assert.Equal(t, "DE*AB7*E840*6487", expectedId.String())
	})
}

func TestEvseId_CompactString(t *testing.T) {
	t.Run("returns a valid ISO string without separators", func(t *testing.T) {
		id := &EvseId{evseid.Id(input.CountryCode, input.OperatorCode, "8406487")}
		assert.Equal(t, "DEAB7E8406487", id.CompactString())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, evseid.Reader, error)
	}{
		{
			name:          "parses a valid ISO EvseId",
			input:         "DE*AB7*E840*6487",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid ISO EvseId without field separators",
			input:         "DEAB7E8406487",
			runAssertions: assertValidCompactId,
		},
		{
			name:          "parses a valid ISO EvseId with an asterisk right after E",
			input:         "DE*AB7*E840*6487",
			runAssertions: assertValidId,
		},
		{
			name:          "returns an error for an ISO EvseId string with an invalid country code",
			input:         "ZZ*AB7*E840*6487",
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error for an invalid ISO EvseId string",
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
		name                                string
		countryCode, operatorCode, outletId string
		checkDigit                          rune
		runAssertions                       func(*testing.T, evseid.Reader, error)
	}{
		{
			name:          "creates an ISO EvseId from its parts",
			countryCode:   input.CountryCode,
			operatorCode:  input.OperatorCode,
			outletId:      input.PowerOutletId,
			runAssertions: assertValidId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			operatorCode:  input.OperatorCode,
			outletId:      input.PowerOutletId,
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			operatorCode:  input.OperatorCode,
			outletId:      input.PowerOutletId,
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   input.CountryCode,
			operatorCode:  "TNMA",
			outletId:      input.PowerOutletId,
			runAssertions: evseid.AssertIsError,
		},
		{
			name:          "returns an error if outletId is too long",
			countryCode:   input.CountryCode,
			operatorCode:  input.OperatorCode,
			outletId:      "0123456789012345678901234567890123456789",
			runAssertions: evseid.AssertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewEvseId(test.countryCode, test.operatorCode, test.outletId)

			test.runAssertions(t, id, err)
		})
	}
}
