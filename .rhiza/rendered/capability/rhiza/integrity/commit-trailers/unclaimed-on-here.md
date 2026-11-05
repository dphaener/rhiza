# unclaimed-on-here

capability · `01M6HM36F6SPZGNYW6GEQKAGZW` · parent [commit-trailers](../commit-trailers.md)

## acceptance

.rhiza/check.json in the change set's commit holds the full hash of a commit on origin/main that is an ancestor of the branch, and make self-check runs the rule and passes

## statement

This repository turns the unclaimed-changes rule on: it commits .rhiza/check.json with unclaimed_since set to the origin/main commit the branch is built on when the file is written

## Links

### Out

- depends-on [unclaimed-commits](unclaimed-commits.md)
- decided-by [unclaimed-changes-cutoff](../../../../decision/unclaimed-changes-cutoff.md)

### In

- verifies [unclaimed-on-here-verified](../../../../evidence/unclaimed-on-here-verified.md)
- verifies [unclaimed-on-here-qa-1](../../../../evidence/unclaimed-on-here-qa-1.md)
- verifies [pr-248-unclaimed-changes](../../../../evidence/pr-248-unclaimed-changes.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
