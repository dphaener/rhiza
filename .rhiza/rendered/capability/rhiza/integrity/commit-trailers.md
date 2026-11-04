# commit-trailers

capability · `01M5SWT0V78KEPDTYECTARQAFS` · parent [integrity](../integrity.md)

## statement

Commits carry a Change trailer that links them to their change set, and a missing one flags an unclaimed change.

## Links

### Out

- decided-by [unclaimed-changes-cutoff](../../../decision/unclaimed-changes-cutoff.md)

### In

- verifies [pr-248-unclaimed-changes](../../../evidence/pr-248-unclaimed-changes.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
