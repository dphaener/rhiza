# plan-unclaimed-changes

source · `01M6HM21KENHTHSEB4PT1E41VV`

## quotes

```json
[]
```

## raw_sha256

b9f389b8e8f05f8118490b95f8b154f847a2e455cc875c72212d7a949101f357

## session

ab09ebc9-684c-463f-93b2-01a1a33c8d88

## summary

Autonomous plan run, no owner. The planner carved leaf commit-trailers out of umbrella Area integrity into this change set and raised one question with three parts: where the unclaimed-changes cutoff lives, which commits count as claimed, and whether this repo turns the rule on now. The arbiter answered with Decision unclaimed-changes-cutoff: the cutoff is a full commit hash in a new committed .rhiza/check.json (repo.json untouched, no file or key means the rule is off, an unresolvable hash is a check error naming the file). The rule reads <cutoff>..HEAD non-merge commits that change a file outside .rhiza/, a commit is claimed when one of its Change: trailers names a change set in the log, and each unclaimed commit gets one warning saying whether the trailer is missing or names an unknown change set. This change set turns the rule on here with the cutoff at the origin/main commit the branch is built on.

## title

Plan of Unclaimed changes, autonomous run

## Links

### In

- derived-from [unclaimed-changes-cutoff](../decision/unclaimed-changes-cutoff.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
