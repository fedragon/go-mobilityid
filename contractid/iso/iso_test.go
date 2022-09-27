package iso

import (
	"github.com/stretchr/testify/assert"
	"mobilityid/contractid"
	"testing"
)

var (
	input = contractid.TestInput{
		CountryCode:   "NL",
		PartyCode:     "TNM",
		InstanceValue: "001234567",
		CheckDigit:    'X',
	}
	expectedId                = &ContractId{contractid.Id(input.CountryCode, input.PartyCode, input.InstanceValue, input.CheckDigit)}
	assertValidId             = contractid.NewAssertValidId(input.CountryCode, input.PartyCode, input.InstanceValue, input.CheckDigit)
	assertValidIdNoCheckDigit = contractid.NewAssertValidIdNoCheckDigit(input.CountryCode, input.PartyCode, input.InstanceValue)
)

func TestContractId_String(t *testing.T) {
	t.Run("returns a valid ISO string with separators", func(t *testing.T) {
		assert.Equal(t, "NL-TNM-001234567-X", expectedId.String())
	})
}

func TestContractId_ToCompactString(t *testing.T) {
	t.Run("returns a valid ISO string with check digit and without separators", func(t *testing.T) {
		assert.Equal(t, "NLTNM001234567X", expectedId.CompactString())
	})
}

func TestContractId_ToCompactStringWithoutCheckDigit(t *testing.T) {
	t.Run("returns a valid ISO string without check digit and separators", func(t *testing.T) {
		assert.Equal(t, "NLTNM001234567", expectedId.CompactStringNoCheckDigit())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, contractid.Reader, error)
	}{
		{
			name:          "parses a valid ISO contract ID with check digit",
			input:         "NL-TNM-001234567-X",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid ISO contract ID without field separators",
			input:         "NLTNM001234567X",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid ISO contract ID without check digit",
			input:         "NL-TNM-001234567",
			runAssertions: assertValidIdNoCheckDigit,
		},
		{
			name:          "parses a valid ISO contract ID without check digit and field separators",
			input:         "NLTNM001234567",
			runAssertions: assertValidIdNoCheckDigit,
		},
		{
			name:          "returns an error for an ISO EvseId string with an invalid country code",
			input:         "ZZ-TNM-001234567-X",
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error for an invalid ISO EvseId string",
			input:         "XYZ",
			runAssertions: contractid.AssertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := Parse(test.input)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewContractId(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		checkDigit                       rune
		runAssertions                    func(*testing.T, contractid.Reader, error)
	}{
		{
			name:          "creates an ISO contract ID from its parts",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			checkDigit:    input.CheckDigit,
			runAssertions: assertValidId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			checkDigit:    input.CheckDigit,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			checkDigit:    input.CheckDigit,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   input.CountryCode,
			partyCode:     "TNMA",
			instance:      input.InstanceValue,
			checkDigit:    input.CheckDigit,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if emi3InstanceValue is too long",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      "001234567890",
			checkDigit:    input.CheckDigit,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if check digit is invalid",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			checkDigit:    'A',
			runAssertions: contractid.AssertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewContractId(test.countryCode, test.partyCode, test.instance, test.checkDigit)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewContractIdNoCheckDigit(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		runAssertions                    func(*testing.T, contractid.Reader, error)
	}{
		{
			name:          "creates an ISO contract ID computing its check digit",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			runAssertions: assertValidId,
		},
		{
			name:          "creates an ISO contract ID from its parts",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			runAssertions: assertValidId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   input.CountryCode,
			partyCode:     "TNMA",
			instance:      input.InstanceValue,
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error if emi3InstanceValue is too long",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      "001234567890",
			runAssertions: contractid.AssertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewContractIdNoCheckDigit(test.countryCode, test.partyCode, test.instance)

			test.runAssertions(t, id, err)
		})
	}
}
