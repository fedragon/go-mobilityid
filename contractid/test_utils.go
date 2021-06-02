package contractid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestInput struct {
	CountryCode string
	PartyCode string
	InstanceValue string
	CheckDigit rune
}

func assertIsError(t *testing.T, _ Reader, err error) {
	assert.NotNil(t, err)
}

func NewAssertValidIdNoCheckDigit(countryCode, partyCode, instance string) func(*testing.T, Reader, error) {
	return func(t *testing.T, reader Reader, err error) {
		assert.Nil(t, err)
		assert.Equal(t, countryCode, reader.CountryCode())
		assert.Equal(t, partyCode, reader.PartyCode())
		assert.Equal(t, instance, reader.InstanceValue())
	}
}

func NewAssertValidId(countryCode, partyCode, instance string, checkDigit rune) func(*testing.T, Reader, error) {
	return func(t *testing.T, reader Reader, err error) {
		assert.Nil(t, err)
		assert.Equal(t, countryCode, reader.CountryCode())
		assert.Equal(t, partyCode, reader.PartyCode())
		assert.Equal(t, instance, reader.InstanceValue())
		assert.Equal(t, checkDigit, reader.CheckDigit())
	}
}
