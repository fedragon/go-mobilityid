package evseid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	isoInput = TestInput{
		CountryCode:   "DE",
		OperatorCode:  "AB7",
		PowerOutletId: "840*6487",
	}
	expectedIsoId           = &IsoEvseId{newEvseId(isoInput.CountryCode, isoInput.OperatorCode, isoInput.PowerOutletId)}
	assertValidIsoId        = NewAssertValidId(isoInput.CountryCode, isoInput.OperatorCode, isoInput.PowerOutletId)
	assertValidIsoCompactId = NewAssertValidCompactId(isoInput.CountryCode, isoInput.OperatorCode, isoInput.PowerOutletId)
)

func TestIsoEvseId_String(t *testing.T) {
	t.Run("returns a valid ISO string with separators", func(t *testing.T) {
		assert.Equal(t, "DE*AB7*E840*6487", expectedIsoId.String())
	})
}

func TestIsoEvseId_CompactString(t *testing.T) {
	t.Run("returns a valid ISO string without separators", func(t *testing.T) {
		id := &IsoEvseId{newEvseId(isoInput.CountryCode, isoInput.OperatorCode, "8406487")}
		assert.Equal(t, "DEAB7E8406487", id.CompactString())
	})
}

func TestParseIso(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, Reader, error)
	}{
		{
			name:          "parses a valid ISO IsoEvseId",
			input:         "DE*AB7*E840*6487",
			runAssertions: assertValidIsoId,
		},
		{
			name:          "parses a valid ISO IsoEvseId without field separators",
			input:         "DEAB7E8406487",
			runAssertions: assertValidIsoCompactId,
		},
		{
			name:          "parses a valid ISO IsoEvseId with an asterisk right after E",
			input:         "DE*AB7*E840*6487",
			runAssertions: assertValidIsoId,
		},
		{
			name:          "returns an error for an ISO IsoEvseId string with an invalid country code",
			input:         "ZZ*AB7*E840*6487",
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error for an invalid ISO IsoEvseId string",
			input:         "XYZ",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := ParseIso(test.input)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewIsoEvseId(t *testing.T) {
	cases := []struct {
		name                                string
		countryCode, operatorCode, outletId string
		checkDigit                          rune
		runAssertions                       func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an ISO IsoEvseId from its parts",
			countryCode:   isoInput.CountryCode,
			operatorCode:  isoInput.OperatorCode,
			outletId:      isoInput.PowerOutletId,
			runAssertions: assertValidIsoId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			operatorCode:  isoInput.OperatorCode,
			outletId:      isoInput.PowerOutletId,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			operatorCode:  isoInput.OperatorCode,
			outletId:      isoInput.PowerOutletId,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   isoInput.CountryCode,
			operatorCode:  "TNMA",
			outletId:      isoInput.PowerOutletId,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if outletId is too long",
			countryCode:   isoInput.CountryCode,
			operatorCode:  isoInput.OperatorCode,
			outletId:      "0123456789012345678901234567890123456789",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewIsoEvseId(test.countryCode, test.operatorCode, test.outletId)

			test.runAssertions(t, id, err)
		})
	}
}
