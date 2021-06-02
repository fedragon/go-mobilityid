package contractid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	isoInput = TestInput{
		CountryCode:   "NL",
		PartyCode:     "TNM",
		InstanceValue: "001234567",
		CheckDigit:    'X',
	}
	isoExpectedId                = &IsoContractId{newContractId(isoInput.CountryCode, isoInput.PartyCode, isoInput.InstanceValue, isoInput.CheckDigit)}
	assertValidIsoId             = NewAssertValidId(isoInput.CountryCode, isoInput.PartyCode, isoInput.InstanceValue, isoInput.CheckDigit)
	assertValidIsoIdNoCheckDigit = NewAssertValidIdNoCheckDigit(isoInput.CountryCode, isoInput.PartyCode, isoInput.InstanceValue)
)

func TestIsoContractId_String(t *testing.T) {
	t.Run("returns a valid ISO string with separators", func(t *testing.T) {
		assert.Equal(t, "NL-TNM-001234567-X", isoExpectedId.String())
	})
}

func TestIsoContractId_ToCompactString(t *testing.T) {
	t.Run("returns a valid ISO string with check digit and without separators", func(t *testing.T) {
		assert.Equal(t, "NLTNM001234567X", isoExpectedId.CompactString())
	})
}

func TestIsoContractId_ToCompactStringWithoutCheckDigit(t *testing.T) {
	t.Run("returns a valid ISO string without check digit and separators", func(t *testing.T) {
		assert.Equal(t, "NLTNM001234567", isoExpectedId.CompactStringNoCheckDigit())
	})
}

func TestParseIso(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, Reader, error)
	}{
		{
			name:          "parses a valid ISO contract ID with check digit",
			input:         "NL-TNM-001234567-X",
			runAssertions: assertValidIsoId,
		},
		{
			name:          "parses a valid ISO contract ID without field separators",
			input:         "NLTNM001234567X",
			runAssertions: assertValidIsoId,
		},
		{
			name:          "parses a valid ISO contract ID without check digit",
			input:         "NL-TNM-001234567",
			runAssertions: assertValidIsoIdNoCheckDigit,
		},
		{
			name:          "parses a valid ISO contract ID without check digit and field separators",
			input:         "NLTNM001234567",
			runAssertions: assertValidIsoIdNoCheckDigit,
		},
		{
			name:          "returns an error for an ISO EvseId string with an invalid country code",
			input:         "ZZ-TNM-001234567-X",
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error for an invalid ISO EvseId string",
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

func TestNewIsoContractId(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		checkDigit                       rune
		runAssertions                    func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an ISO contract ID from its parts",
			countryCode:   isoInput.CountryCode,
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			checkDigit:    isoInput.CheckDigit,
			runAssertions: assertValidIsoId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			checkDigit:    isoInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			checkDigit:    isoInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   isoInput.CountryCode,
			partyCode:     "TNMA",
			instance:      isoInput.InstanceValue,
			checkDigit:    isoInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if emi3InstanceValue is too long",
			countryCode:   isoInput.CountryCode,
			partyCode:     isoInput.PartyCode,
			instance:      "001234567890",
			checkDigit:    isoInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if check digit is invalid",
			countryCode:   isoInput.CountryCode,
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			checkDigit:    'A',
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewIsoContractId(test.countryCode, test.partyCode, test.instance, test.checkDigit)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewIsoContractIdNoCheckDigit(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		runAssertions                    func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an ISO contract ID computing its check digit",
			countryCode:   isoInput.CountryCode,
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			runAssertions: assertValidIsoId,
		},
		{
			name:          "creates an ISO contract ID from its parts",
			countryCode:   isoInput.CountryCode,
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			runAssertions: assertValidIsoId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     isoInput.PartyCode,
			instance:      isoInput.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   isoInput.CountryCode,
			partyCode:     "TNMA",
			instance:      isoInput.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if emi3InstanceValue is too long",
			countryCode:   isoInput.CountryCode,
			partyCode:     isoInput.PartyCode,
			instance:      "001234567890",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewIsoContractIdNoCheckDigit(test.countryCode, test.partyCode, test.instance)

			test.runAssertions(t, id, err)
		})
	}
}
