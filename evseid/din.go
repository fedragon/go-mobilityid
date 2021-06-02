package evseid

import (
	"fmt"
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

	if i := len(code); i < 1 || i > 3 {
		return fmt.Errorf("country code '%s' has leading '+' but doesn't match expected length (2 to 4)", code)
	}

	for _, r := range code {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("country code '%s' can only contain a leading (optional) '+' following by 1 to 3 digits", countryCode)
		}
	}

	if i := len(operatorCode); i < 3 || i > 6 {
		return fmt.Errorf("operator code '%s' doesn't match length (2 < x < 7)", operatorCode)
	}

	if i := len(powerOutletId); i == 0 || i > 31 {
		return fmt.Errorf("power outlet id '%s' doesn't match expected length (0 < x < 32)", powerOutletId)
	}

	return nil
}
