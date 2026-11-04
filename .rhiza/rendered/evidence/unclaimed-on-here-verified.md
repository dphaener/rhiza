# unclaimed-on-here-verified

evidence · `01M6KEA3NCVH4DJZ5K9MTNVMHD` · actor [claude-code-3e50350b](../actor/claude-code-3e50350b.md)

## result

PASS

## statement

Verifier claude-code-3e50350b at leaf commit f60fd4ba (attempt 2). .rhiza/check.json holds unclaimed_since=18e33781c268cf557f5309e13d2402706ac67881: a full 40-char hash, git cat-file says commit, merge-base --is-ancestor confirms it is an ancestor of HEAD and of origin/main, and it equals merge-base HEAD origin/main. With RHIZA_ROOT/RHIZA_CHANGE/RHIZA_RUN unset, make self-check exits 0 (ok, base 18e33781); the attempt-1 finding (rendered/actor/claude-code-d222929a.md with no tracked record) is fixed: that file and its index line are removed, and the newly rendered claude-code-9fe165a9 has tracked log records at HEAD. The rule is on: internal/api/check.go unclaimedCutoff reads .rhiza/check.json, rhiza check --json at HEAD has no findings, and in a scratch clone at f60fd4ba with one extra commit lacking a Change: trailer, rhiza check --json reports rule=unclaimed severity=warn for that commit.

## Links

### Out

- verifies [unclaimed-on-here](../capability/rhiza/integrity/commit-trailers/unclaimed-on-here.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
