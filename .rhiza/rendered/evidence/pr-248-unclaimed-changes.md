# pr-248-unclaimed-changes

evidence · `01M6KPQYPZ1WNJ6QAQWQT41XZV` · actor [dphaener](../actor/dphaener.md)

## head

b82582b04db296de471599b63493531dbb202e31

## pr

https://github.com/dphaener/rhiza/pull/248

## result

1 check passed: check SUCCESS

## statement

Pull request #248 (https://github.com/dphaener/rhiza/pull/248) passed CI at b82582b

## status

verified

## Links

### Out

- verifies [commit-trailers](../capability/rhiza/integrity/commit-trailers.md)
- verifies [unclaimed-cutoff](../capability/rhiza/integrity/commit-trailers/unclaimed-cutoff.md)
- verifies [unclaimed-commits](../capability/rhiza/integrity/commit-trailers/unclaimed-commits.md)
- verifies [unclaimed-on-here](../capability/rhiza/integrity/commit-trailers/unclaimed-on-here.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
