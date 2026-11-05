# unclaimed-cutoff-qa-1

evidence · `01M6KERKPMWMAK21SHJ14YY8WG` · actor [claude-code-dc452748](../actor/claude-code-dc452748.md)

## result

PASS

## statement

QA at 0bee2696 with a scratch build (go build -o $TMPDIR/rhiza-x ./cmd/rhiza) on a scratch git repo with rhiza init, a real change set and untrailered code commits. rhiza check with no .rhiza/check.json and with check.json {} printed ok, no unclaimed line, exit 0. unclaimed_since 0123456789012345678901234567890123456789 exited 1 with 'rhiza: .rhiza/check.json: unclaimed_since "0123..." is not a commit in this git repository'; '-h' and malformed JSON also exited 1 naming .rhiza/check.json. repo.json was not touched.

## Links

### Out

- verifies [unclaimed-cutoff](../capability/rhiza/integrity/commit-trailers/unclaimed-cutoff.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
