# estimatelaws

Go package that embeds the dataset "U.S. State Auto Repair Estimate Laws": one record per U.S. state and the District of Columbia (51 rows), release 2026-09-30. No dependency outside the standard library.

**Not legal advice. Verify the official text.** This dataset summarizes public legal texts for research and consumer information. Laws change, and a short summary can miss exceptions. Read the official text, or ask a lawyer, before relying on any row.

## Install

    go get github.com/repairworth/estimatelaws

## Use

```go
law, ok := estimatelaws.Lookup("CA")
if ok && law.Status == estimatelaws.Verified {
    fmt.Println(law.StatuteCitation, law.OverrunWithoutNewConsent)
}
for _, l := range estimatelaws.All() {
    fmt.Println(l.Code, l.Status)
}
fmt.Println(estimatelaws.Version)
```

## Status values

| Status | Meaning | Rows |
|---|---|---|
| `verified` | The official text was read and the row names a statute or rule. Rule values come from that text. | 33 |
| `no_statute_found` | An official source was read on the date given and no dedicated statute was found (`official_url`, `checked_on`). This is a negative finding: it cannot be proven by quotation. Rule cells are empty. | 2 |
| `lead` | The official text was not read; any citation or reference kept in the record is a pointer to follow up, and the rule fields are empty. | 16 |

An empty field means "not established from the source read". It does not mean "no rule". Only `verified` rows carry rule values.

## Source

- State-by-state table and estimate checker: https://repairworth.com/car-repair-estimate/
- Data sources: https://repairworth.com/data-sources/
- Publisher: https://repairworth.com/

The embedded file is `data/us-state-auto-repair-estimate-laws.csv`. Its SHA-256 is recorded in `data/SHA256SUMS` and checked by the tests.

## License

Code: MIT, see `LICENSE`. Data in `data/`: Creative Commons Attribution 4.0 International, see `data/LICENSE` (https://creativecommons.org/licenses/by/4.0/). The statutes and regulations quoted are government edicts in the public domain.
