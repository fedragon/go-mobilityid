package evseid

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
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

	if err := v.Validate(operatorCode, v.Required, v.Length(3, 3)); err != nil {
		return err
	}

	if err := v.Validate(powerOutletId, v.Required, v.Length(1, 31)); err != nil {
		return err
	}

	return nil
}
