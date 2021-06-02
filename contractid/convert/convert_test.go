package convert

import (
	"github.com/stretchr/testify/assert"
	"mobilityid/contractid"
	"testing"
)

func TestDinToEmi3(t *testing.T) {
	t.Run("converts a DIN into a valid EMI3 contract id if it has a leading 0", func(t *testing.T) {
		dinContractId, err := contractid.NewDinContractId(
			"NL",
			"TNM",
			"012204",
			'5',
		)
		assert.Nil(t, err)

		emi3ContractId, err := DinToEmi3(dinContractId)

		assert.Nil(t, err)
		assert.Equal(t, "NL-TNM-C00122045-K", emi3ContractId.String())
	})
}

func TestEmi3ToDin(t *testing.T) {
	t.Run("converts an EMI3 into a valid DIN contract id", func(t *testing.T) {
		emi3ContractId, err := contractid.NewEmi3ContractId(
			"NL",
			"TNM",
			"00122045",
			'K',
		)
		assert.Nil(t, err)

		dinContractId, err := Emi3ToDin(emi3ContractId)

		assert.Nil(t, err)
		assert.Equal(t, "NL-TNM-012204-5", dinContractId.String())
	})

	t.Run("returns an error if it does not have a leading 0", func(t *testing.T) {
		emi3ContractId, err := contractid.NewEmi3ContractId(
			"NL",
			"TNM",
			"33122045",
			'P',
		)
		assert.Nil(t, err)

		_, err = Emi3ToDin(emi3ContractId)

		assert.NotNil(t, err)
	})
}

func TestEmi3ToIso(t *testing.T) {
	t.Run("converts an EMI3 into a valid ISO contract id", func(t *testing.T) {
		emi3ContractId, err := contractid.NewEmi3ContractId(
			"NL",
			"TNM",
			"00122045",
			'K',
		)
		assert.Nil(t, err)

		isoContractId, err := Emi3ToIso(emi3ContractId)

		assert.Nil(t, err)
		assert.Equal(t, "NL-TNM-C00122045-K", isoContractId.String())
	})
}
