package evseid

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type TestInput struct {
	CountryCode string
	OperatorCode string
	PowerOutletId string
}
func assertIsError(t *testing.T, _ Reader, err error) {
	assert.NotNil(t, err)
}

func NewAssertValidId(countryCode, operatorCode, powerOutletId string) func(*testing.T, Reader, error) {
	return func(t *testing.T, reader Reader, err error) {
		assert.Nil(t, err)
		assert.Equal(t, countryCode, reader.CountryCode())
		assert.Equal(t, operatorCode, reader.OperatorCode())
		assert.Equal(t, powerOutletId, reader.PowerOutletId())
	}
}

func NewAssertValidCompactId(countryCode, operatorCode, powerOutletId string) func(*testing.T, Reader, error) {
	return func(t *testing.T, reader Reader, err error) {
		assert.Nil(t, err)
		assert.Equal(t, countryCode, reader.CountryCode())
		assert.Equal(t, operatorCode, reader.OperatorCode())
		assert.Equal(t, strings.ReplaceAll(powerOutletId, "*", ""), reader.PowerOutletId())
	}
}
