package estimatelaws

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
)

const csvName = "us-state-auto-repair-estimate-laws.csv"

// The embedded file must be the exported file, byte for byte.
func TestEmbeddedChecksum(t *testing.T) {
	sums, err := os.ReadFile("data/SHA256SUMS")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(rawCSV)
	got := hex.EncodeToString(sum[:])
	for _, line := range strings.Split(string(sums), "\n") {
		if f := strings.Fields(line); len(f) == 2 && f[1] == csvName {
			if f[0] != got {
				t.Fatalf("sha256 of embedded %s = %s, SHA256SUMS says %s", csvName, got, f[0])
			}
			return
		}
	}
	t.Fatalf("%s is not listed in data/SHA256SUMS", csvName)
}

func TestAll(t *testing.T) {
	all := All()
	if len(all) != 51 {
		t.Fatalf("All() returned %d records, want 51", len(all))
	}
	seen := map[string]bool{}
	count := map[Status]int{}
	latest := ""
	for _, l := range all {
		count[l.Status]++
		if len(l.Code) != 2 || l.State == "" || seen[l.Code] {
			t.Errorf("bad or duplicate record: %q %q", l.Code, l.State)
		}
		seen[l.Code] = true
		if l.CheckedOn > latest {
			latest = l.CheckedOn
		}
		rules := l.WrittenEstimateRequired + l.OverrunWithoutNewConsent + l.ConsentForm +
			l.ReturnOfReplacedParts + l.ShopRegistrationRequired + l.DiagnosticOrTeardownFeeDisclosure
		switch l.Status {
		case Verified:
			if l.StatuteCitation == "" || l.OfficialURL == "" || (l.HasStatute != "yes" && l.HasStatute != "partial") {
				t.Errorf("%s: verified record without citation, official URL or statute flag", l.Code)
			}
		case NoStatuteFound:
			if rules != "" || l.HasThreshold || l.StatuteCitation != "" || l.HasStatute != "no" || l.OfficialURL == "" || l.Note == "" {
				t.Errorf("%s: no_statute_found record must carry no rule value, an official URL and a note", l.Code)
			}
		case Lead:
			if rules != "" || l.HasThreshold || l.HasStatute != "" || l.Note == "" {
				t.Errorf("%s: lead record must carry no rule value, no statute flag and a note", l.Code)
			}
		default:
			t.Errorf("%s: unknown status %q", l.Code, l.Status)
		}
	}
	want := map[Status]int{Verified: 33, NoStatuteFound: 2, Lead: 16}
	for status, n := range want {
		if count[status] != n {
			t.Errorf("%d %s records, want %d", count[status], status, n)
		}
	}
	if latest != Version {
		t.Errorf("Version = %q, most recent CheckedOn = %q", Version, latest)
	}
	all[0].State = "changed"
	if All()[0].State == "changed" {
		t.Error("All() must return a copy")
	}
}

func TestLookup(t *testing.T) {
	for _, code := range []string{"DC", "dc", " Dc "} {
		if l, ok := Lookup(code); !ok || l.State != "District of Columbia" {
			t.Errorf("Lookup(%q) = %q, %v", code, l.State, ok)
		}
	}
	for _, code := range []string{"", "ZZ", "California"} {
		if _, ok := Lookup(code); ok {
			t.Errorf("Lookup(%q) should not match", code)
		}
	}
}

func ExampleLookup() {
	law, ok := Lookup("ca")
	if !ok {
		return
	}
	fmt.Println(law.State)
	fmt.Println(law.Status)
	fmt.Println(law.StatuteCitation)
	// Output:
	// California
	// verified
	// Cal. Bus. & Prof. Code §9884.9
}

func ExampleAll() {
	fmt.Println(len(All()), "jurisdictions, release", Version)
	// Output:
	// 51 jurisdictions, release 2026-09-30
}
