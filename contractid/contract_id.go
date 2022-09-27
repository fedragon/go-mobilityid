package contractid

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
	c "mobilityid/common"
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
	err := v.Validate(
		countryCode,
		v.Required,
		v.Length(2, 2),
		v.By(
			func(value interface{}) error {
				if !c.IsValidCountryCode(value.(string)) {
					return fmt.Errorf("country code '%s' is not valid", value.(string))
				}
				return nil
			}),
	)
	if err != nil {
		return err
	}

	if err := v.Validate(partyCode, v.Required, v.Length(3, 3)); err != nil {
		return err
	}

	if err := v.Validate(instance, v.Required, v.Length(1, instanceMaxLength)); err != nil {
		return err
	}

	return nil
}
