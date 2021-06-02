package convert

import (
	"errors"
	"fmt"
	"mobilityid/contractid"
	"strings"
)

// DinToEmi3 converts a DIN contract ID to its EMI3 equivalent
func DinToEmi3(id *contractid.DinContractId) (*contractid.Emi3ContractId, error) {
	return contractid.NewEmi3ContractIdNoCheckDigit(id.CountryCode(), id.PartyCode(), fmt.Sprintf("0%s%c", id.InstanceValue(), id.CheckDigit()))
}

// Emi3ToDin converts an EMI3 contract ID to its DIN equivalent, if possible
func Emi3ToDin(id *contractid.Emi3ContractId) (*contractid.DinContractId, error) {
	if !strings.HasPrefix(id.InstanceValue(), "0") {
		return nil, errors.New("cannot convert to DIN: instance value is too long")
	}

	dinInstance := id.InstanceValue()[1:7]
	dinCheckDigit := id.InstanceValue()[7:8][0]

	return contractid.NewDinContractId(id.CountryCode(), id.PartyCode(), dinInstance, rune(dinCheckDigit))
}

// Emi3ToIso converts an EMI3 contract ID to its ISO equivalent, if possible
func Emi3ToIso(id *contractid.Emi3ContractId) (*contractid.IsoContractId, error) {
	return contractid.NewIsoContractId(id.CountryCode(), id.PartyCode(), fmt.Sprintf("C%s", id.InstanceValue()), id.CheckDigit())
}
