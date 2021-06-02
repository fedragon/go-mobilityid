package contractid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	dinInput = TestInput{
		CountryCode:   "IN",
		PartyCode:     "TNM",
		InstanceValue: "000071",
		CheckDigit:    '9',
	}
	dinExpectedId                 = &DinContractId{newContractId(dinInput.CountryCode, dinInput.PartyCode, dinInput.InstanceValue, dinInput.CheckDigit)}
	assertValidDinId             = NewAssertValidId(dinInput.CountryCode, dinInput.PartyCode, dinInput.InstanceValue, dinInput.CheckDigit)
	assertValidDinIdNoCheckDigit = NewAssertValidIdNoCheckDigit(dinInput.CountryCode, dinInput.PartyCode, dinInput.InstanceValue)
)

func TestDinContractId_String(t *testing.T) {
	t.Run("returns a valid DIN string with separators", func(t *testing.T) {
		assert.Equal(t, "IN-TNM-000071-9", dinExpectedId.String())
	})
}

func TestDinContractId_ToCompactString(t *testing.T) {
	t.Run("returns a valid DIN string with check digit and without separators", func(t *testing.T) {
		assert.Equal(t, "INTNM0000719", dinExpectedId.CompactString())
	})
}

func TestDinContractId_ToCompactStringWithoutCheckDigit(t *testing.T) {
	t.Run("returns a valid DIN string without check digit and separators", func(t *testing.T) {
		assert.Equal(t, "INTNM000071", dinExpectedId.CompactStringNoCheckDigit())
	})
}

func TestParse(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, Reader, error)
	}{
		{
			name:          "parses a valid DIN contract ID with check digit",
			input:         "IN-TNM-000071-9",
			runAssertions: assertValidDinId,
		},
		{
			name:          "parses a valid DIN contract ID without field separators",
			input:         "INTNM0000719",
			runAssertions: assertValidDinId,
		},
		{
			name:          "parses a valid DIN contract ID without check digit",
			input:         "IN-TNM-000071",
			runAssertions: assertValidDinIdNoCheckDigit,
		},
		{
			name:          "parses a valid DIN contract ID without check digit and field separators",
			input:         "INTNM000071",
			runAssertions: assertValidDinIdNoCheckDigit,
		},
		{
			name:          "returns an error for an DIN EvseId string with an invalid country code",
			input:         "ZZ-TNM-000071-9",
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error for an invalid DIN EvseId string",
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

func TestNewContractId(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		checkDigit                       rune
		runAssertions                    func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an DIN contract ID from its parts",
			countryCode:   dinInput.CountryCode,
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			checkDigit:    dinInput.CheckDigit,
			runAssertions: assertValidDinId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			checkDigit:    dinInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			checkDigit:    dinInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   dinInput.CountryCode,
			partyCode:     "TNMA",
			instance:      dinInput.InstanceValue,
			checkDigit:    dinInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if instance is too long",
			countryCode:   dinInput.CountryCode,
			partyCode:     dinInput.PartyCode,
			instance:      "001234567890",
			checkDigit:    dinInput.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if check digit is invalid",
			countryCode:   dinInput.CountryCode,
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			checkDigit:    'A',
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewDinContractId(test.countryCode, test.partyCode, test.instance, test.checkDigit)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewContractIdNoCheckDigit(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		runAssertions                    func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an DIN contract ID computing its check digit",
			countryCode:   dinInput.CountryCode,
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			runAssertions: assertValidDinId,
		},
		{
			name:          "creates an DIN contract ID from its parts",
			countryCode:   dinInput.CountryCode,
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			runAssertions: assertValidDinId,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     dinInput.PartyCode,
			instance:      dinInput.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   dinInput.CountryCode,
			partyCode:     "TNMA",
			instance:      dinInput.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if instance is too long",
			countryCode:   dinInput.CountryCode,
			partyCode:     dinInput.PartyCode,
			instance:      "001234567890",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewDinContractIdNoCheckDigit(test.countryCode, test.partyCode, test.instance)

			test.runAssertions(t, id, err)
		})
	}
}
