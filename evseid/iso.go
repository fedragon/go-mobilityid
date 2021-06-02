package evseid

import (
	"fmt"
	c "mobilityid/common"
	"regexp"
	"strings"
)

var isoRegex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)(?:\\*?)(?P<operator>%v)(?:\\*?)%v(?P<outlet>%v)$", c.CountryCodeRegex, c.PartyCodeRegex, "[Ee]", "([A-Za-z0-9\\*]{1,31})"))

// IsoEvseId represents an ISO EVSE ID
type IsoEvseId struct {
	Reader
}

func (c IsoEvseId) String() string {
	return strings.Join([]string{c.CountryCode(), c.OperatorCode(), "E" + c.PowerOutletId()}, "*")
}

func (c IsoEvseId) CompactString() string {
	return strings.ReplaceAll(c.String(), "*", "")
}

// NewIsoEvseId returns an IsoEvseId, if provided input is valid; returns an error otherwise.
func NewIsoEvseId(countryCode, operatorCode, powerOutletId string) (*IsoEvseId, error) {
	if err := validateIso(countryCode, operatorCode, powerOutletId); err != nil {
		return nil, err
	}

	return &IsoEvseId{
		newEvseId(
			strings.ToUpper(countryCode),
			strings.ToUpper(operatorCode),
			strings.ToUpper(powerOutletId),
		),
	}, nil
}

// ParseIso parses the input string into an IsoEvseId, if it is valid; returns an error otherwise.
func ParseIso(input string) (*IsoEvseId, error) {
	groups := isoRegex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO IsoEvseId: %v", input)
	}

	if !c.IsValidCountryCode(countryCode) {
		return nil, fmt.Errorf("invalid country code: %v", countryCode)
	}

	operatorCode, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "operator", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO IsoEvseId: %v", input)
	}
	outletId, err := c.ExtractAndUpcaseGroup(isoRegex, groups, "outlet", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO IsoEvseId: %v", input)
	}

	return &IsoEvseId{
		newEvseId(
			countryCode,
			operatorCode,
			outletId,
		),
	}, nil
}

func validateIso(countryCode, operatorCode, powerOutletId string) error {
	if len(countryCode) != 2 {
		return fmt.Errorf("country code '%s' doesn't match expected length (2)", countryCode)
	}

	if !c.IsValidCountryCode(countryCode) {
		return fmt.Errorf("country code '%s' is not valid", countryCode)
	}

	if len(operatorCode) != 3 {
		return fmt.Errorf("operator code '%s' doesn't match length (3)", operatorCode)
	}

	if i := len(powerOutletId); i == 0 || i > 31 {
		return fmt.Errorf("power outlet id '%s' doesn't match expected length (0 < x < 32)", powerOutletId)
	}

	return nil
}
