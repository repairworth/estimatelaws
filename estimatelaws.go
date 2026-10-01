// Package estimatelaws embeds the dataset "U.S. State Auto Repair Estimate
// Laws": one record for each of the 50 states and the District of Columbia,
// describing the state law or regulation on written auto repair estimates.
//
// This is not legal advice. Verify the official text before relying on any
// record: laws change, and a short summary can miss exceptions. A Verified
// record carries the citation and the URL of the source to read.
//
// Only records whose Status is Verified carry rule values. They were read on
// an official state source. For a NoStatuteFound record,
// an official source was read on the date given and no dedicated statute was
// found (OfficialURL, CheckedOn); an absence cannot be proven by quotation.
// For a Lead record,
// the official text was not read; any citation or reference kept in the
// record is a pointer to follow up, and the rule fields are empty. An empty
// field means "not established from the source read"; it never means "no
// rule".
//
// The dataset is compiled and published by RepairWorth:
//
//   - state-by-state table: https://repairworth.com/car-repair-estimate/
//   - data sources: https://repairworth.com/data-sources/
//   - publisher: https://repairworth.com/
//
// The code is under the MIT license. The embedded data, in the data
// directory, is under Creative Commons Attribution 4.0 International.
package estimatelaws

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// Version is the release date of the embedded dataset: the most recent
// CheckedOn date found in it.
const Version = "2026-09-30"

//go:embed data/us-state-auto-repair-estimate-laws.csv
var rawCSV []byte

// Status says how far a record can be relied on.
type Status string

const (
	// Verified: the official text was read and the record names a statute or rule.
	Verified Status = "verified"
	// NoStatuteFound: an official source was read on the date given and no
	// dedicated statute was found.
	NoStatuteFound Status = "no_statute_found"
	// Lead: the official text was not read; any citation or reference kept
	// in the record is a pointer to follow up, and the rule fields are empty.
	Lead Status = "lead"
)

// Law is one row of the dataset. Text fields hold the wording of the dataset;
// an empty string means "not established from the source read".
type Law struct {
	State  string // jurisdiction name
	Code   string // two-letter USPS code
	Status Status

	HasStatute      string // "yes" or "partial" (Verified), "no" (NoStatuteFound), "" (Lead)
	StatuteCitation string
	OfficialURL     string
	SecondaryURL    string // legal mirror, never used as evidence

	WrittenEstimateRequired string
	// WrittenEstimateThresholdUSD is meaningful only when HasThreshold is true.
	WrittenEstimateThresholdUSD       int
	HasThreshold                      bool
	OverrunWithoutNewConsent          string
	ConsentForm                       string
	ReturnOfReplacedParts             string
	ShopRegistrationRequired          string
	DiagnosticOrTeardownFeeDisclosure string

	VerifiedQuote    string // verbatim quote from the official source
	CheckedOn        string // ISO 8601 date of the reading or of the search
	EffectiveVersion string
	Note             string
}

var laws, byCode = mustLoad(rawCSV)

// Lookup returns the record for a two-letter USPS code such as "CA" or "DC".
// The code is matched without regard to case or surrounding spaces.
func Lookup(code string) (Law, bool) {
	l, ok := byCode[strings.ToUpper(strings.TrimSpace(code))]
	return l, ok
}

// All returns every record, in the order of the dataset. The slice is a
// copy: the caller may modify it.
func All() []Law {
	return append([]Law(nil), laws...)
}

func mustLoad(raw []byte) ([]Law, map[string]Law) {
	records, err := csv.NewReader(bytes.NewReader(raw)).ReadAll()
	if err != nil || len(records) < 2 {
		panic(fmt.Sprintf("estimatelaws: embedded CSV is unreadable: %v", err))
	}
	col := make(map[string]int, len(records[0]))
	for i, name := range records[0] {
		col[name] = i
	}
	get := func(rec []string, name string) string {
		i, ok := col[name]
		if !ok {
			panic("estimatelaws: embedded CSV lacks column " + name)
		}
		return rec[i]
	}
	list := make([]Law, 0, len(records)-1)
	index := make(map[string]Law, len(records)-1)
	for _, rec := range records[1:] {
		l := Law{
			State:                             get(rec, "state"),
			Code:                              get(rec, "state_code"),
			Status:                            Status(get(rec, "status")),
			HasStatute:                        get(rec, "has_statute"),
			StatuteCitation:                   get(rec, "statute_citation"),
			OfficialURL:                       get(rec, "official_url"),
			SecondaryURL:                      get(rec, "secondary_url"),
			WrittenEstimateRequired:           get(rec, "written_estimate_required"),
			OverrunWithoutNewConsent:          get(rec, "overrun_without_new_consent"),
			ConsentForm:                       get(rec, "consent_form"),
			ReturnOfReplacedParts:             get(rec, "return_of_replaced_parts"),
			ShopRegistrationRequired:          get(rec, "shop_registration_required"),
			DiagnosticOrTeardownFeeDisclosure: get(rec, "diagnostic_or_teardown_fee_disclosure"),
			VerifiedQuote:                     get(rec, "verified_quote"),
			CheckedOn:                         get(rec, "checked_on"),
			EffectiveVersion:                  get(rec, "effective_version"),
			Note:                              get(rec, "note"),
		}
		if t := get(rec, "written_estimate_threshold_usd"); t != "" {
			n, err := strconv.Atoi(t)
			if err != nil {
				panic("estimatelaws: threshold is not an integer for " + l.Code)
			}
			l.WrittenEstimateThresholdUSD, l.HasThreshold = n, true
		}
		list = append(list, l)
		index[l.Code] = l
	}
	return list, index
}
