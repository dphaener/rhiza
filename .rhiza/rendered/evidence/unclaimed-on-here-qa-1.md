# unclaimed-on-here-qa-1

evidence · `01M6KET8BK20TKYJF932Y22FA0` · actor [claude-code-dc452748](../actor/claude-code-dc452748.md)

## result

PASS

## statement

QA at 0bee2696: .rhiza/check.json holds unclaimed_since 18e33781c268cf557f5309e13d2402706ac67881, a commit (git cat-file), an ancestor of HEAD and of origin/main, equal to merge-base HEAD origin/main. With RHIZA_* unset, make self-check exits 0 (ok) and scratch-build rhiza check --json has no findings. In a clone of this checkout with an added untrailered code commit, rhiza check reported 'warn unclaimed: commit ... "stray code" ... has no Change: trailer (missing)', so the rule is on here; an added empty commit raised nothing.

## Links

### Out

- verifies [unclaimed-on-here](../capability/rhiza/integrity/commit-trailers/unclaimed-on-here.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
