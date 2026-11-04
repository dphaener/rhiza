# unclaimed-cutoff

capability · `01M6HM2WJCT93N5GDJ35P10P4B` · parent [commit-trailers](../commit-trailers.md)

## acceptance

Testscripts run rhiza check with no check.json, with check.json lacking unclaimed_since, and with an unresolvable hash: the first two report no unclaimed finding and exit 0 on an otherwise clean repo with untrailered commits, the third fails with a message naming .rhiza/check.json

## statement

rhiza check reads the unclaimed-changes cutoff from an optional committed .rhiza/check.json, {"unclaimed_since": "<full commit hash>"}: no file or no key turns the rule off, and a hash that does not resolve to a commit is a check error naming the file; repo.json is untouched

## Links

### Out

- decided-by [unclaimed-changes-cutoff](../../../../decision/unclaimed-changes-cutoff.md)

### In

- depends-on [unclaimed-commits](unclaimed-commits.md)
- verifies [unclaimed-cutoff-verified](../../../../evidence/unclaimed-cutoff-verified.md)
- verifies [unclaimed-cutoff-qa-1](../../../../evidence/unclaimed-cutoff-qa-1.md)
- verifies [pr-248-unclaimed-changes](../../../../evidence/pr-248-unclaimed-changes.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
