package contractid

import (
	"fmt"
	c "mobilityid/common"
	"mobilityid/contractid/din"
	"regexp"
	"strings"
)

var dinRegex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)(?:[*-]?)(?P<party>%v)(?:[*-]?)(?P<instance>%v)(?:(?:[*-]?)(?P<check>%v))?$", c.CountryCodeRegex, c.PartyCodeRegex, "([A-Za-z0-9]{6})", c.CheckDigitRegex))

type DinContractId struct {
	Reader
}

// NewDinContractIdNoCheckDigit returns a DIN contract ID complete of check digit, if provided input is valid; returns an error otherwise.
func NewDinContractIdNoCheckDigit(countryCode, partyCode, instance string) (*DinContractId, error) {
	if err := validateNoCheckDigit(countryCode, partyCode, instance, 6); err != nil {
		return nil, err
	}

	return &DinContractId{
		newContractId(
			strings.ToUpper(countryCode),
			strings.ToUpper(partyCode),
			strings.ToUpper(instance),
			din.ComputeCheckDigit(countryCode+partyCode+instance),
		),
	}, nil
}

// NewDinContractId returns a DIN contract ID, if provided input is valid; returns an error otherwise.
func NewDinContractId(countryCode, partyCode, instance string, checkDigit rune) (*DinContractId, error) {
	id, err := NewDinContractIdNoCheckDigit(countryCode, partyCode, instance)
	if err != nil {
		return nil, err
	}

	if checkDigit != id.CheckDigit() {
		return nil, fmt.Errorf("provided check digit '%c' doesn't match computed one '%c'", checkDigit, id.CheckDigit())
	}

	return id, nil
}

// ParseDin parses the input string into a DIN contract ID, if it is valid; returns an error otherwise.
// A check digit will only be present, in returned struct, if the provided string contained it.
func ParseDin(input string) (*DinContractId, error) {
	groups := dinRegex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN contract ID: %v", input)
	}
	partyCode, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "party", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN contract ID: %v", input)
	}
	instance, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "instance", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN contract ID: %v", input)
	}
	check, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "check", false)
	if err != nil {
		return nil, fmt.Errorf("not a DIN contract ID: %v", input)
	}

	var checkDigit rune
	if len(check) > 0 {
		checkDigit = rune(check[0])
		if err := validateDin(countryCode, partyCode, instance, checkDigit); err != nil {
			return nil, err
		}
	} else if err := validateNoCheckDigit(countryCode, partyCode, instance, 6); err != nil {
		return nil, err
	}

	return &DinContractId{
		newContractId(
			countryCode,
			partyCode,
			instance,
			din.ComputeCheckDigit(countryCode+partyCode+instance),
		),
	}, nil
}

func validateDin(countryCode, partyCode, instance string, checkDigit rune) error {
	if err := validateNoCheckDigit(countryCode, partyCode, instance, 6); err != nil {
		return err
	}

	if checkDigit != din.ComputeCheckDigit(strings.ToUpper(countryCode+partyCode+instance)) {
		return fmt.Errorf("check digit '%c' is invalid", checkDigit)
	}

	return nil
}
