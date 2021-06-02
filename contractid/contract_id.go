package contractid

import (
	"errors"
	"fmt"
	"mobilityid/common"
	"strings"
)

// Stringer provides functions to get string representations of contract IDs
type Stringer interface {
	String() string
	CompactString() string
	CompactStringNoCheckDigit() string
}

// Reader provides functions to read fields of contract IDs
type Reader interface {
	CountryCode() string
	PartyCode() string
	InstanceValue() string
	CheckDigit() rune
	PartyId() string
	CompactPartyId() string
	Stringer
}

type contractId struct {
	countryCode   string
	partyCode     string
	instanceValue string
	checkDigit    rune
	Reader
}

func newContractId(countryCode, partyCode, instanceValue string, checkDigit rune) *contractId {
	return &contractId{
		countryCode:   countryCode,
		partyCode:     partyCode,
		instanceValue: instanceValue,
		checkDigit:    checkDigit,
	}
}

// CountryCode returns the contract's country code
func (id *contractId) CountryCode() string {
	return id.countryCode
}

// PartyCode returns the contract's party code
func (id *contractId) PartyCode() string {
	return id.partyCode
}

// InstanceValue returns the contract's emi3InstanceValue value
func (id *contractId) InstanceValue() string {
	return id.instanceValue
}

// CheckDigit returns the contract's check digit
func (id *contractId) CheckDigit() rune {
	return id.checkDigit
}

// PartyId returns this contracts' party ID
func (id *contractId) PartyId() string {
	return id.CountryCode() + "-" + id.PartyCode()
}

// CompactPartyId returns this contracts' party ID without separator
func (id *contractId) CompactPartyId() string {
	return id.CountryCode() + id.PartyCode()
}

// String returns a canonical contract ID string representation
func (id *contractId) String() string {
	result := fmt.Sprintf("%s-%s-%s", id.CountryCode(), id.PartyCode(), id.InstanceValue())
	if id.CheckDigit() != '0' {
		result = fmt.Sprintf("%s-%c", result, id.CheckDigit())
	}

	return result
}

// CompactString returns a contract ID string without separators
func (id *contractId) CompactString() string {
	return strings.ReplaceAll(id.String(), "-", "")
}

// CompactStringNoCheckDigit returns a contract ID string without separators nor check digit
func (id *contractId) CompactStringNoCheckDigit() string {
	compact := id.CompactString()
	return compact[:len(compact)-1]
}

// validateNoCheckDigit validates provided inputs
func validateNoCheckDigit(countryCode, partyCode, instance string, instanceMaxLength int) error {
	if len(countryCode) != 2 {
		return fmt.Errorf("country code '%s' doesn't match expected length (2)", countryCode)
	}

	if !common.IsValidCountryCode(countryCode) {
		return fmt.Errorf("country code '%s' is not valid", countryCode)
	}

	if len(partyCode) != 3 {
		return fmt.Errorf("party code '%s' doesn't match length (3)", partyCode)
	}

	if len(instance) == 0 {
		return errors.New("instance value cannot be empty")
	}

	if len(instance) > instanceMaxLength {
		return fmt.Errorf("instance value '%s' exceeds max length (%v)", instance, instanceMaxLength)
	}

	return nil
}
