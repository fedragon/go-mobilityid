package iso

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
	c "mobilityid/common"
	"mobilityid/evseid"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)(?:\\*?)(?P<operator>%v)(?:\\*?)%v(?P<outlet>%v)$", c.CountryCodeRegex, c.PartyCodeRegex, "[Ee]", "([A-Za-z0-9\\*]{1,31})"))

// EvseId represents an ISO EVSE ID
type EvseId struct {
	evseid.Reader
}

func (c EvseId) String() string {
	return strings.Join([]string{c.CountryCode(), c.OperatorCode(), "E" + c.PowerOutletId()}, "*")
}

func (c EvseId) CompactString() string {
	return strings.ReplaceAll(c.String(), "*", "")
}

// NewEvseId returns an EvseId, if provided input is valid; returns an error otherwise.
func NewEvseId(countryCode, operatorCode, powerOutletId string) (*EvseId, error) {
	if err := validate(countryCode, operatorCode, powerOutletId); err != nil {
		return nil, err
	}

	return &EvseId{
		evseid.Id(
			strings.ToUpper(countryCode),
			strings.ToUpper(operatorCode),
			strings.ToUpper(powerOutletId),
		),
	}, nil
}

// Parse parses the input string into an EvseId, if it is valid; returns an error otherwise.
func Parse(input string) (*EvseId, error) {
	groups := regex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(regex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO EvseId: %v", input)
	}

	if !c.IsValidCountryCode(countryCode) {
		return nil, fmt.Errorf("invalid country code: %v", countryCode)
	}

	operatorCode, err := c.ExtractAndUpcaseGroup(regex, groups, "operator", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO EvseId: %v", input)
	}
	outletId, err := c.ExtractAndUpcaseGroup(regex, groups, "outlet", true)
	if err != nil {
		return nil, fmt.Errorf("not an ISO EvseId: %v", input)
	}

	return &EvseId{
		evseid.Id(
			countryCode,
			operatorCode,
			outletId,
		),
	}, nil
}

func validate(countryCode, operatorCode, powerOutletId string) error {
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
