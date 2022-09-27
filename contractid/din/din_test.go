package din

import (
	"github.com/stretchr/testify/assert"
	"mobilityid/contractid"
	"testing"
)

var (
	input = contractid.TestInput{
		CountryCode:   "IN",
		PartyCode:     "TNM",
		InstanceValue: "000071",
		CheckDigit:    '9',
	}
	expectedId                = &ContractId{contractid.Id(input.CountryCode, input.PartyCode, input.InstanceValue, input.CheckDigit)}
	assertValidId             = contractid.NewAssertValidId(input.CountryCode, input.PartyCode, input.InstanceValue, input.CheckDigit)
	assertValidIdNoCheckDigit = contractid.NewAssertValidIdNoCheckDigit(input.CountryCode, input.PartyCode, input.InstanceValue)
)

func TestContractId_String(t *testing.T) {
	t.Run("returns a valid DIN string with separators", func(t *testing.T) {
		assert.Equal(t, "IN-TNM-000071-9", expectedId.String())
	})
}

func TestContractId_ToCompactString(t *testing.T) {
	t.Run("returns a valid DIN string with check digit and without separators", func(t *testing.T) {
		assert.Equal(t, "INTNM0000719", expectedId.CompactString())
	})
}

func TestContractId_ToCompactStringWithoutCheckDigit(t *testing.T) {
	t.Run("returns a valid DIN string without check digit and separators", func(t *testing.T) {
		assert.Equal(t, "INTNM000071", expectedId.CompactStringNoCheckDigit())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, contractid.Reader, error)
	}{
		{
			name:          "parses a valid DIN contract ID with check digit",
			input:         "IN-TNM-000071-9",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid DIN contract ID without field separators",
			input:         "INTNM0000719",
			runAssertions: assertValidId,
		},
		{
			name:          "parses a valid DIN contract ID without check digit",
			input:         "IN-TNM-000071",
			runAssertions: assertValidIdNoCheckDigit,
		},
		{
			name:          "parses a valid DIN contract ID without check digit and field separators",
			input:         "INTNM000071",
			runAssertions: assertValidIdNoCheckDigit,
		},
		{
			name:          "returns an error for an DIN EvseId string with an invalid country code",
			input:         "ZZ-TNM-000071-9",
			runAssertions: contractid.AssertIsError,
		},
		{
			name:          "returns an error for an invalid DIN EvseId string",
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
			name:          "creates an DIN contract ID from its parts",
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
			name:          "returns an error if instance is too long",
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
			name:          "creates an DIN contract ID computing its check digit",
			countryCode:   input.CountryCode,
			partyCode:     input.PartyCode,
			instance:      input.InstanceValue,
			runAssertions: assertValidId,
		},
		{
			name:          "creates an DIN contract ID from its parts",
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
			name:          "returns an error if instance is too long",
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
