# unclaimed-cutoff-verified

evidence · `01M6K22GD7PWE84X7NG5JWSQ2W` · actor [claude-code-5d92770f](../actor/claude-code-5d92770f.md)

## result

PASS

## statement

At 50b2a773 ran go test ./cmd/rhiza -run TestScripts/check-unclaimed-cutoff -count=1 -v: in a repo with an untrailered code commit, rhiza check --json with no .rhiza/check.json and with check.json {} printed "ok": true and no unclaimed text (exit 0); a resolving unclaimed_since printed ok; unclaimed_since 0123456789012345678901234567890123456789 exited 1 with stderr 'rhiza: .rhiza/check.json: unclaimed_since "0123..." is not a commit in this git repository'. Read the diff: internal/api/check.go unclaimedCutoff reads store.Dir/check.json, missing file or key gives "" (rule off), resolves via git rev-parse --verify <h>^{commit}, rejects empty and dash-prefixed values, and passes the hash to check.Input.UnclaimedSince; repo.json is not touched. make test: every package ok (internal/check, internal/api, cmd/rhiza included); the TestReadBudget timing budget afterwards failed under heavy machine load (api package took 739s), in the read path this diff does not touch.

## Links

### Out

- verifies [unclaimed-cutoff](../capability/rhiza/integrity/commit-trailers/unclaimed-cutoff.md)

## History

- `01M6HKRVC93TGX0ZC40QZ2GW9J` Unclaimed changes: rhiza check warns about commits since a configured cutoff that change files outside .rhiza/ with no Change: trailer (internal/check, internal/api/check.go parseTrailerLog). Relates to leaf commit-trailers in umbrella 01M5SWT0K9TV1R (Area integrity); carve it if it still stands there
