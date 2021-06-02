package contractid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	emi3Input = TestInput{
		CountryCode:   "NL",
		PartyCode:     "TNM",
		InstanceValue: "00122045",
		CheckDigit:    'K',
	}
	emi3ExpectedId                = &Emi3ContractId{newContractId(emi3Input.CountryCode, emi3Input.PartyCode, emi3Input.InstanceValue, emi3Input.CheckDigit)}
	assertValidEmi3Id             = NewAssertValidId(emi3Input.CountryCode, emi3Input.PartyCode, emi3Input.InstanceValue, emi3Input.CheckDigit)
	assertValidEmi3IdNoCheckDigit = NewAssertValidIdNoCheckDigit(emi3Input.CountryCode, emi3Input.PartyCode, emi3Input.InstanceValue)
)

func TestEmi3ContractId_String(t *testing.T) {
	t.Run("returns a valid EMI3 string with separators", func(t *testing.T) {
		assert.Equal(t, "NL-TNM-C00122045-K", emi3ExpectedId.String())
	})
}

func TestEmi3ContractId_CompactString(t *testing.T) {
	t.Run("returns a valid EMI3 string with check digit and without separators", func(t *testing.T) {
		assert.Equal(t, "NLTNMC00122045K", emi3ExpectedId.CompactString())
	})
}

func TestEmi3ContractId_CompactStringWithoutCheckDigit(t *testing.T) {
	t.Run("returns a valid EMI3 string without check digit and separators", func(t *testing.T) {
		assert.Equal(t, "NLTNMC00122045", emi3ExpectedId.CompactStringNoCheckDigit())
	})
}

func TestParseEmi3(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		runAssertions func(*testing.T, Reader, error)
	}{
		{
			name:          "parses a valid EMI3 contract ID with check digit",
			input:         "NL-TNM-C00122045-K",
			runAssertions: assertValidEmi3Id,
		},
		{
			name:          "parses a valid EMI3 contract ID without field separators",
			input:         "NLTNMC00122045K",
			runAssertions: assertValidEmi3Id,
		},
		{
			name:          "parses a valid EMI3 contract ID with mixed case",
			input:         "NltNMc00122045k",
			runAssertions: assertValidEmi3Id,
		},
		{
			name:          "parses a valid EMI3 contract ID without check digit",
			input:         "NL-TNM-C00122045",
			runAssertions: assertValidEmi3IdNoCheckDigit,
		},
		{
			name:          "parses a valid EMI3 contract ID without check digit and field separators",
			input:         "NLTNMC00122045",
			runAssertions: assertValidEmi3IdNoCheckDigit,
		},
		{
			name:          "returns an error for an EMI3 EvseId string with an invalid country code",
			input:         "ZZ-TNM-C00122045",
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error for an invalid EMI3 EvseId string",
			input:         "XYZ",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := ParseEmi3(test.input)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewEmi3ContractId(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		checkDigit                       rune
		runAssertions                    func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an EMI3 contract ID from its parts",
			countryCode:   emi3Input.CountryCode,
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			checkDigit:    emi3Input.CheckDigit,
			runAssertions: assertValidEmi3Id,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			checkDigit:    emi3Input.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			checkDigit:    emi3Input.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   emi3Input.CountryCode,
			partyCode:     "TNMA",
			instance:      emi3Input.InstanceValue,
			checkDigit:    emi3Input.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if emi3Input.InstanceValue is too long",
			countryCode:   emi3Input.CountryCode,
			partyCode:     emi3Input.PartyCode,
			instance:      "C001234567890",
			checkDigit:    emi3Input.CheckDigit,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if check digit is invalid",
			countryCode:   emi3Input.CountryCode,
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			checkDigit:    'A',
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewEmi3ContractId(test.countryCode, test.partyCode, test.instance, test.checkDigit)

			test.runAssertions(t, id, err)
		})
	}
}

func TestNewEmi3ContractIdNoCheckDigit(t *testing.T) {
	cases := []struct {
		name                             string
		countryCode, partyCode, instance string
		runAssertions                    func(*testing.T, Reader, error)
	}{
		{
			name:          "creates an EMI3 contract ID computing its check digit",
			countryCode:   emi3Input.CountryCode,
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			runAssertions: assertValidEmi3Id,
		},
		{
			name:          "creates an EMI3 contract ID from its parts",
			countryCode:   emi3Input.CountryCode,
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			runAssertions: assertValidEmi3Id,
		},
		{
			name:          "returns an error if invalid country code",
			countryCode:   "ZZ",
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if country code is too long",
			countryCode:   "XYZ",
			partyCode:     emi3Input.PartyCode,
			instance:      emi3Input.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if party code is too long",
			countryCode:   emi3Input.CountryCode,
			partyCode:     "TNMA",
			instance:      emi3Input.InstanceValue,
			runAssertions: assertIsError,
		},
		{
			name:          "returns an error if emi3Input.InstanceValue is too long",
			countryCode:   emi3Input.CountryCode,
			partyCode:     emi3Input.PartyCode,
			instance:      "001234567890",
			runAssertions: assertIsError,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewEmi3ContractIdNoCheckDigit(test.countryCode, test.partyCode, test.instance)

			test.runAssertions(t, id, err)
		})
	}
}
