package iso

import (
	"fmt"
	c "mobilityid/common"
	"mobilityid/contractid"
	"regexp"
	"strings"
)

const instanceMaxLength = 9

var regex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)(?:-?)(?P<party>%v)(?:-?)(?P<emi3InstanceValue>%v)(?:(?:-?)(?P<check>%v))?$", c.CountryCodeRegex, c.PartyCodeRegex, "([A-Za-z0-9]{9})", c.CheckDigitRegex))

// ContractId represents an ISO15118-1 contract identifier
type ContractId struct {
	contractid.Reader
}

// NewContractIdNoCheckDigit returns an ISO contract ID complete of check digit, if provided input is valid; returns an error otherwise.
func NewContractIdNoCheckDigit(countryCode, partyCode, instance string) (*ContractId, error) {
	if err := contractid.ValidateNoCheckDigit(countryCode, partyCode, instance, instanceMaxLength); err != nil {
		return nil, err
	}

	checkDigit, err := ComputeCheckDigit(strings.ToUpper(countryCode + partyCode + instance))
	if err != nil {
		return nil, fmt.Errorf("unable to compute check digit: %w", err)
	}

	return &ContractId{
		contractid.Id(
			strings.ToUpper(countryCode),
			strings.ToUpper(partyCode),
			strings.ToUpper(instance),
			checkDigit,
		),
	}, nil
}

// NewContractId returns an ISO contract ID, if provided input is valid; returns an error otherwise.
func NewContractId(countryCode, partyCode, instance string, checkDigit rune) (*ContractId, error) {
	id, err := NewContractIdNoCheckDigit(countryCode, partyCode, instance)
	if err != nil {
		return nil, err
	}

	if checkDigit != id.CheckDigit() {
		return nil, fmt.Errorf("provided check digit '%c' doesn't match computed one '%c'", checkDigit, id.CheckDigit())
	}

	return id, nil
}

// Parse parses the input string into an ISO contract ID, if it is valid; returns an error otherwise.
// A check digit will only be present, in returned struct, if the provided string contained it.
func Parse(input string) (*ContractId, error) {
	groups := regex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(regex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}
	partyCode, err := c.ExtractAndUpcaseGroup(regex, groups, "party", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}
	instance, err := c.ExtractAndUpcaseGroup(regex, groups, "emi3InstanceValue", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}
	check, err := c.ExtractAndUpcaseGroup(regex, groups, "check", false)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}

	var checkDigit rune
	if len(check) > 0 {
		checkDigit = rune(check[0])
		if err := validate(countryCode, partyCode, instance, checkDigit); err != nil {
			return nil, err
		}
	} else if err := contractid.ValidateNoCheckDigit(countryCode, partyCode, instance, instanceMaxLength); err != nil {
		return nil, err
	}

	return &ContractId{
		contractid.Id(
			countryCode,
			partyCode,
			instance,
			checkDigit,
		),
	}, nil
}

func validate(countryCode, partyCode, instance string, checkDigit rune) error {
	if err := contractid.ValidateNoCheckDigit(countryCode, partyCode, instance, instanceMaxLength); err != nil {
		return err
	}

	computed, err := ComputeCheckDigit(countryCode + partyCode + instance)
	if err != nil {
		return fmt.Errorf("unable to compute check digit: %w", err)
	}

	if checkDigit != computed {
		return fmt.Errorf("check digit '%c' is invalid", checkDigit)
	}

	return nil
}
