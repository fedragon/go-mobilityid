package evseid

// Reader provides functions to read fields of evse IDs
type Reader interface {
	CountryCode() string
	OperatorCode() string
	PowerOutletId() string
	PartyId() string
	CompactPartyId() string
}

type evseId struct {
	countryCode   string
	operatorCode  string
	powerOutletId string
	Reader
}

func Id(countryCode, partyCode, powerOutletId string) *evseId {
	return &evseId{
		countryCode:   countryCode,
		operatorCode:  partyCode,
		powerOutletId: powerOutletId,
	}
}

// CountryCode returns the country code
func (id *evseId) CountryCode() string {
	return id.countryCode
}

// OperatorCode returns the party code
func (id *evseId) OperatorCode() string {
	return id.operatorCode
}

// PowerOutletId returns the power outlet ID
func (id *evseId) PowerOutletId() string {
	return id.powerOutletId
}

// PartyId returns the party ID
func (id *evseId) PartyId() string {
	return id.CountryCode() + "-" + id.OperatorCode()
}

// CompactPartyId returns the party ID without separator
func (id *evseId) CompactPartyId() string {
	return id.CountryCode() + id.OperatorCode()
}
