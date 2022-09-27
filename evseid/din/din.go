package din

import (
	"fmt"
	v "github.com/go-ozzo/ozzo-validation"
	c "mobilityid/common"
	"mobilityid/evseid"
	"regexp"
	"strings"
	"unicode"
)

var regex = regexp.MustCompile(fmt.Sprintf("^(?P<country>%v)\\*(?P<operator>%v)\\*(?P<outlet>%v)$", "\\+?([0-9]{1,3})", "([0-9]{3,6})", "([0-9\\*]{1,32})"))

// EvseId represents a DIN Evse Id
type EvseId struct {
	evseid.Reader
}

func (c EvseId) String() string {
	result := strings.Join([]string{c.CountryCode(), c.OperatorCode(), c.PowerOutletId()}, "*")

	return result
}

// NewEvseId returns a DIN EvseId, if provided input is valid; returns an error otherwise.
func NewEvseId(countryCode, operatorCode, powerOutletId string) (*EvseId, error) {
	if err := validate(countryCode, operatorCode, powerOutletId); err != nil {
		return nil, err
	}

	return Parse(strings.Join([]string{countryCode, operatorCode, powerOutletId}, "*"))
}

// Parse parses the input string into a EvseId, if it is valid; returns an error otherwise.
func Parse(input string) (*EvseId, error) {
	groups := regex.FindStringSubmatch(input)

	countryCode, err := c.ExtractAndUpcaseGroup(regex, groups, "country", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN EvseId: %v", input)
	}
	if !strings.HasPrefix(countryCode, "+") {
		countryCode = "+" + countryCode
	}
	operatorId, err := c.ExtractAndUpcaseGroup(regex, groups, "operator", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN EvseId: %v", input)
	}
	outletId, err := c.ExtractAndUpcaseGroup(regex, groups, "outlet", true)
	if err != nil {
		return nil, fmt.Errorf("not a DIN EvseId: %v", input)
	}

	return &EvseId{
		evseid.Id(
			countryCode,
			operatorId,
			outletId,
		),
	}, nil
}

func validate(countryCode, operatorCode, powerOutletId string) error {
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
