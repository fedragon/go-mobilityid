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

func newEvseId(countryCode, partyCode, powerOutletId string) *evseId {
	return &evseId{
		countryCode:   countryCode,
		operatorCode:  partyCode,
		powerOutletId: powerOutletId,
	}
}

// CountryCode returns the evse's country code
func (id *evseId) CountryCode() string {
	return id.countryCode
}

// OperatorCode returns the evse's party code
func (id *evseId) OperatorCode() string {
	return id.operatorCode
}

// PowerOutletId returns the evse's emi3InstanceValue value
func (id *evseId) PowerOutletId() string {
	return id.powerOutletId
}

// PartyId returns this evses' party ID
func (id *evseId) PartyId() string {
	return id.CountryCode() + "-" + id.OperatorCode()
}

// CompactPartyId returns this evses' party ID without separator
func (id *evseId) CompactPartyId() string {
	return id.CountryCode() + id.OperatorCode()
}
