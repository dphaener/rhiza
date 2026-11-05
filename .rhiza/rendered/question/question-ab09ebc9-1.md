# question-ab09ebc9-1

question · `01M6HKY1F0HDWE1BP7DSZVM3MF`

## session

ab09ebc9-684c-463f-93b2-01a1a33c8d88

## statement

Unclaimed changes: where does the cutoff (commits after which rhiza check warns on a missing Change: trailer) live? It must be committed so CI sees it; repo.json is currently a strict canonical format-1 file with format/id/remotes. / Unclaimed changes: which commits after the cutoff count as claimed? (Only commits that change a file outside .rhiza/ are considered, as the evidence-actor rule already does.) / Unclaimed changes: should this repo turn the rule on as part of this change set, with the cutoff set to the main commit it lands on?

## Links

### In

- answers [unclaimed-changes-cutoff](../decision/unclaimed-changes-cutoff.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
