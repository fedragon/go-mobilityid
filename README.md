# go-mobilityid

Porting of the Scala [MobilityId](https://github.com/ShellRechargeSolutionsEU/mobilityid/) library to Go.

## Features

- Create instances of DIN91826, ISO15118-1, or eMI3 contract IDs
- Compute (or validate, if provided) their check digit

## Usage

```go
// Contract IDs

emi3Id, err := contractid.NewEmi3ContractIdNoCheckDigit("NL", "TNM", "00122045")
if err != nil {
  // ...
}

emi3Id, err = contractid.NewEmi3ContractId("NL", "TNM", "00122045", 'K')
if err != nil {
  // ...
}

fmt.Println(emi3Id.CountryCode()) // "NL"
fmt.Println(emi3Id.PartyCode()) // "TNM"
fmt.Println(emi3Id.InstanceValue()) // "00122045"
fmt.Println(emi3Id.CheckDigit()) // 'K'

fmt.Println(emi3Id.PartyId()) // "NL-TNM"
fmt.Println(emi3Id.CompactPartyId()) // "NLTNM"

fmt.Println(emi3Id.String()) // "NL-TNM-C00122045-K"
fmt.Println(emi3Id.CompactString()) // "NLTNMC00122045K"
fmt.Println(emi3Id.CompactStringNoCheckDigit()) // "NLTNMC00122045"

dinId, err := convert.Emi3ToDin(id)
if err != nil {
  // ...
}

fmt.Println(dinId.String()) // "NL-TNM-012204-5"

// EVSE IDs

isoId, err := evseid.NewIsoEvseId("NL", "TNM", "030123456*0")
if err != nil {
  // ...
}

fmt.Println(isoId.CountryCode()) // "NL"
fmt.Println(isoId.OperatorCode()) // "TNM"
fmt.Println(isoId.PowerOutletId()) // "030123456*0"

fmt.Println(isoId.String()) // "NL*TNM*E030123456*0"
fmt.Println(isoId.CompactString()) // "NLTNME0301234560"
```

## Differences with original library

### EMI3 instance value

Given the EMI3 contract ID `NL-TNM-C12345678-J`, getting the `instanceValue`:

- from the Scala implementation (at the moment of writing, `v1.1.0`), will return `C12345678`
- from this library's implementation, will return `12345678`

This is because the leading `C` character is considered part of the format, rather than of the _instance value_.

### ContractId conversions

Direct conversions from `ISO` to `DIN` and vice versa, which are deprecated in the original library, have not been
ported.

## Caveats

Types are public so that they can be referenced by clients, which means that it _is_ technically possible to initialize
them directly (e.g. `contractid.Emi3ContractId{}`); this will, however, create an empty instance which cannot be modified in any
way because it only provides functions exported by the `contractId.Reader` interface, so it's effectively of no use to a
client.
