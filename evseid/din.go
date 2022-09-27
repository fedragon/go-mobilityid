package evseid

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
	c "mobilityid/common"
	"regexp"
	"strings"
	"unicode"
)

var dinRegex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)\\*(?P<operator>%v)\\*(?P<outlet>%v)$", "\\+?([0-9]{1,3})", "([0-9]{3,6})", "([0-9\\*]{1,32})"))

// DinEvseId represents a DIN Evse Id
type DinEvseId struct {
	Reader
}

func (c DinEvseId) String() string {
	result := strings.Join([]string{c.CountryCode(), c.OperatorCode(), c.PowerOutletId()}, "*")

	return result
}

// NewDinEvseId returns a DinEvseId, if provided input is valid; returns an error otherwise.
func NewDinEvseId(countryCode, operatorCode, powerOutletId string) (*DinEvseId, error) {
	if err := validateDin(countryCode, operatorCode, powerOutletId); err != nil {
		return nil, err
	}

	return ParseDin(strings.Join([]string{countryCode, operatorCode, powerOutletId}, "*"))
}

// ParseDin parses the input string into a DinEvseId, if it is valid; returns an error otherwise.
func ParseDin(input string) (*DinEvseId, error) {
	groups := dinRegex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN DinEvseId: %v", input)
	}
	if !strings.HasPrefix(countryCode, "+") {
		countryCode = "+" + countryCode
	}
	operatorId, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "operator", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN DinEvseId: %v", input)
	}
	outletId, err := c.ExtractAndUpcaseGroup(dinRegex, groups, "outlet", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN DinEvseId: %v", input)
	}

	return &DinEvseId{
		newEvseId(
			strings.ToUpper(countryCode),
			strings.ToUpper(operatorId),
			strings.ToUpper(outletId),
		),
	}, nil
}

func validateDin(countryCode, operatorCode, powerOutletId string) error {
	code := countryCode
	if strings.HasPrefix(countryCode, "+") {
		code = countryCode[1:]
	}

	if err := v.Validate(code, v.Required, v.Length(2, 4)); err != nil {
		return err
	}

	err := v.Validate(code, v.By(func(value interface{}) error {
		num := value.(string)
		for _, r := range num {
			if !unicode.IsDigit(r) {
				return fmt.Errorf("country code '%s' can only contain a leading (optional) '+' following by 1 to 3 digits", countryCode)
			}
		}

		return nil
	}))
	if err != nil {
		return err
	}

	if err := v.Validate(operatorCode, v.Required, v.Length(2, 6)); err != nil {
		return err
	}

	if err := v.Validate(powerOutletId, v.Required, v.Length(1, 31)); err != nil {
		return err
	}

	return nil
}
