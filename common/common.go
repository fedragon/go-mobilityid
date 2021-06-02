package common

import (
	"fmt"
	"github.com/pariz/gountries"
	"regexp"
	"strings"
)

const (
	CountryCodeRegex = "([A-Za-z]{2})"
	PartyCodeRegex   = "([A-Za-z0-9]{3})"
	CheckDigitRegex  = "([A-Za-z0-9])"
)

var countries = gountries.New()

// ExtractAndUpcaseGroup extracts a group (if found) and returns its value in upper case
func ExtractAndUpcaseGroup(re *regexp.Regexp, groups []string, name string, required bool) (string, error) {
	if index := re.SubexpIndex(name); index < len(groups) {
		return strings.ToUpper(groups[index]), nil
	}

	if !required {
		return "", nil
	}

	return "", fmt.Errorf("group not found: %v", name)
}

// IsValidCountryCode validates if the country code exists
func IsValidCountryCode(code string) bool {
	_, err := countries.FindCountryByAlpha(code)

	return err == nil
}

