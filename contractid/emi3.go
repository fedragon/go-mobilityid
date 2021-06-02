package contractid

import (
	"fmt"
	c "mobilityid/common"
	"mobilityid/contractid/iso"
	"regexp"
	"strings"
)

var emi3Regex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)(?:-?)(?P<party>%v)(?:-?)%s(?P<emi3InstanceValue>%v)(?:(?:-?)(?P<check>%v))?$", c.CountryCodeRegex, c.PartyCodeRegex, "[Cc]", "([A-Za-z0-9]{8})", c.CheckDigitRegex))

// Emi3ContractId represents an EMI3 contract identifier
type Emi3ContractId struct {
	Reader
}

// Overriding some Stringer interface methods in order to include 'C' id type

func (id *Emi3ContractId) String() string {
	result := fmt.Sprintf("%s-%s-%c%s", id.CountryCode(), id.PartyCode(), 'C', id.InstanceValue())
	if id.CheckDigit() != '0' {
		result = fmt.Sprintf("%s-%c", result, id.CheckDigit())
	}

	return result
}

func (id *Emi3ContractId) CompactString() string {
	return strings.ReplaceAll(id.String(), "-", "")
}

func (id *Emi3ContractId) CompactStringNoCheckDigit() string {
	compact := id.CompactString()
	return compact[:len(compact)-1]
}

// NewEmi3ContractIdNoCheckDigit returns an EMI3 contract ID complete of check digit, if provided input is valid; returns an error otherwise.
func NewEmi3ContractIdNoCheckDigit(countryCode, partyCode, instance string) (*Emi3ContractId, error) {
	if err := validateNoCheckDigit(countryCode, partyCode, instance, 8); err != nil {
		return nil, err
	}

	code := strings.ToUpper(fmt.Sprintf("%s%s%c%s", countryCode, partyCode, 'C', instance))
	checkDigit, err := iso.ComputeCheckDigit(code)
	if err != nil {
		return nil, fmt.Errorf("unable to compute check digit: %w", err)
	}

	return &Emi3ContractId{
		newContractId(
			strings.ToUpper(countryCode),
			strings.ToUpper(partyCode),
			strings.ToUpper(instance),
			checkDigit,
		),
	}, nil
}

// NewEmi3ContractId returns an EMI3 contract ID, if provided input is valid; returns an error otherwise.
func NewEmi3ContractId(countryCode, partyCode, instance string, checkDigit rune) (*Emi3ContractId, error) {
	id, err := NewEmi3ContractIdNoCheckDigit(countryCode, partyCode, instance)
	if err != nil {
		return nil, err
	}

	if checkDigit != id.CheckDigit() {
		return nil, fmt.Errorf("provided check digit ('%v') doesn't match computed one ('%v')", checkDigit, id.CheckDigit())
	}

	return id, nil
}

// ParseEmi3 parses the input string into an EMI3 contract ID, if it is valid; returns an error otherwise.
// A check digit will only be present, in returned struct, if the provided string contained it.
func ParseEmi3(input string) (*Emi3ContractId, error) {
	groups := emi3Regex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(emi3Regex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not an EMI3 contract ID: %v", input)
	}
	partyCode, err := c.ExtractAndUpcaseGroup(emi3Regex, groups, "party", true)
	if err != nil {
		return nil, fmt.Errorf("not an EMI3 contract ID: %v", input)
	}
	instance, err := c.ExtractAndUpcaseGroup(emi3Regex, groups, "emi3InstanceValue", true)
	if err != nil {
		return nil, fmt.Errorf("not an EMI3 contract ID: %v", input)
	}
	check, err := c.ExtractAndUpcaseGroup(emi3Regex, groups, "check", false)
	if err != nil {
		return nil, fmt.Errorf("not an EMI3 contract ID: %v", input)
	}

	var checkDigit rune
	if len(check) > 0 {
		checkDigit = rune(check[0])
		if err := validateEmi3(countryCode, partyCode, instance, checkDigit); err != nil {
			return nil, err
		}
	} else if err := validateNoCheckDigit(countryCode, partyCode, instance, 8); err != nil {
		return nil, err
	}

	return &Emi3ContractId{
		newContractId(
			strings.ToUpper(countryCode),
			strings.ToUpper(partyCode),
			strings.ToUpper(instance),
			checkDigit,
		),
	}, nil
}

func validateEmi3(countryCode, partyCode, instance string, checkDigit rune) error {
	if err := validateNoCheckDigit(countryCode, partyCode, instance, 8); err != nil {
		return err
	}

	code := strings.ToUpper(fmt.Sprintf("%s%s%c%s", countryCode, partyCode, 'C', instance))
	computed, err := iso.ComputeCheckDigit(code)
	if err != nil {
		return fmt.Errorf("unable to compute check digit: %w", err)
	}

	if checkDigit != computed {
		return fmt.Errorf("check digit '%c' is invalid", checkDigit)
	}

	return nil
}
