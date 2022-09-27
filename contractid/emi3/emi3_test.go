package emi3

import (
	"github.com/stretchr/testify/assert"
	"mobilityid/contractid"
	"testing"
)

var (
	input = contractid.TestInput{
		CountryCode:   "NL",
		PartyCode:     "TNM",
		InstanceValue: "00122045",
		CheckDigit:    'K',
	}
	expectedId                = &ContractId{contractid.Id(input.CountryCode, input.PartyCode, input.InstanceValue, input.CheckDigit)}
	assertValidId             = contractid.NewAssertValidId(input.CountryCode, input.PartyCode, input.InstanceValue, input.CheckDigit)
	assertValidIdNoCheckDigit = contractid.NewAssertValidIdNoCheckDigit(input.CountryCode, input.PartyCode, input.InstanceValue)
)

func TestContractId_String(t *testing.T) {
	t.Run("returns a valid EMI3 string with separators", func(t *testing.T) {
		assert.Equal(t, "NL-TNM-C00122045-K", expectedId.String())
	})
}

func TestContractId_CompactString(t *testing.T) {
	t.Run("returns a valid EMI3 string with check digit and without separators", func(t *testing.T) {
		assert.Equal(t, "NLTNMC00122045K", expectedId.CompactString())
	})
}

func TestContractId_CompactStringWithoutCheckDigit(t *testing.T) {
	t.Run("returns a valid EMI3 string without check digit and separators", func(t *testing.T) {
		assert.Equal(t, "NLTNMC00122045", expectedId.CompactStringNoCheckDigit())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, contractid.Reader, error)
	}{
		{
			name:          "parses a valid EMI3 contract ID with check digit",
			input:         "NL-TNM-C00122045-K",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid EMI3 contract ID without field separators",
			input:         "NLTNMC00122045K",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid EMI3 contract ID with mixed case",
			input:         "NltNMc00122045k",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid EMI3 contract ID without check digit",
			input:         "NL-TNM-C00122045",
			runAssertions: assertValidIdNoCheckDigit,
		},
		{
			name:          "parses a valid EMI3 contract ID without check digit and field separators",
			input:         "NLTNMC00122045",
			runAssertions: assertValidIdNoCheckDigit,
		},
		{
			name:          "returns an error for an EMI3 EvseId string with an invalid country code",
			input:         "ZZ-TNM-C00122045",
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error for an invalid EMI3 EvseId string",
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
			name:          "creates an EMI3 contract ID from its parts",
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
			name:          "returns an error if input.InstanceValue is too long",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      "C001234567890",
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
			name:          "creates an EMI3 contract ID computing its check digit",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			runAssertions: assertValidId,
		},
		{
			name:          "creates an EMI3 contract ID from its parts",
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
			name:          "returns an error if input.InstanceValue is too long",
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
