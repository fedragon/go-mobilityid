package contractid

import (
	"fmt"
	c "mobilityid/common"
	"mobilityid/contractid/iso"
	"regexp"
	"strings"
)

var isoRegex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)(?:-?)(?P<party>%v)(?:-?)(?P<emi3InstanceValue>%v)(?:(?:-?)(?P<check>%v))?$", c.CountryCodeRegex, c.PartyCodeRegex, "([A-Za-z0-9]{9})", c.CheckDigitRegex))

// IsoContractId represents an ISO15118-1 contract identifier
type IsoContractId struct {
	Reader
}

// NewIsoContractIdNoCheckDigit returns an ISO contract ID complete of check digit, if provided input is valid; returns an error otherwise.
func NewIsoContractIdNoCheckDigit(countryCode, partyCode, instance string) (*IsoContractId, error) {
	if err := validateNoCheckDigit(countryCode, partyCode, instance, 9); err != nil {
		return nil, err
	}

	checkDigit, err := iso.ComputeCheckDigit(strings.ToUpper(countryCode + partyCode + instance))
	if err != nil {
		return nil, fmt.Errorf("unable to compute check digit: %w", err)
	}

	return &IsoContractId{
		newContractId(
			strings.ToUpper(countryCode),
			strings.ToUpper(partyCode),
			strings.ToUpper(instance),
			checkDigit,
		),
	}, nil
}

// NewIsoContractId returns an ISO contract ID, if provided input is valid; returns an error otherwise.
func NewIsoContractId(countryCode, partyCode, instance string, checkDigit rune) (*IsoContractId, error) {
	id, err := NewIsoContractIdNoCheckDigit(countryCode, partyCode, instance)
	if err != nil {
		return nil, err
	}

	if checkDigit != id.CheckDigit() {
		return nil, fmt.Errorf("provided check digit '%c' doesn't match computed one '%c'", checkDigit, id.CheckDigit())
	}

	return id, nil
}

// ParseIso parses the input string into an ISO contract ID, if it is valid; returns an error otherwise.
// A check digit will only be present, in returned struct, if the provided string contained it.
func ParseIso(input string) (*IsoContractId, error) {
	groups := isoRegex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}
	partyCode, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "party", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}
	instance, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "emi3InstanceValue", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}
	check, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "check", false)
	if err != nil {
		return nil, fmt.Errorf("not an ISO contract ID: %v", input)
	}

	var checkDigit rune
	if len(check) > 0 {
		checkDigit = rune(check[0])
		if err := validate(countryCode, partyCode, instance, checkDigit); err != nil {
			return nil, err
		}
	} else if err := validateNoCheckDigit(countryCode, partyCode, instance, 9); err != nil {
		return nil, err
	}

	return &IsoContractId{
		newContractId(
			strings.ToUpper(countryCode),
			strings.ToUpper(partyCode),
			strings.ToUpper(instance),
			checkDigit,
		),
	}, nil
}

func validate(countryCode, partyCode, instance string, checkDigit rune) error {
	if err := validateNoCheckDigit(countryCode, partyCode, instance, 9); err != nil {
		return err
	}

	computed, err := iso.ComputeCheckDigit(countryCode + partyCode + instance)
	if err != nil {
		return fmt.Errorf("unable to compute check digit: %w", err)
	}

	if checkDigit != computed {
		return fmt.Errorf("check digit '%c' is invalid", checkDigit)
	}

	return nil
}
