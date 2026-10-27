# build-check-verified

evidence · `01M5ZMT7QCQ3NZFDY7N0RHAR87` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Testscript build-check: with Evidence for alpha by the verifier actor and none for beta, build-check exits 1 naming beta ('no Evidence verifies it'); with beta's Evidence by the builder, whose commit carries its Actor: and Leaf: beta trailers, it exits 1 naming beta, the Evidence, the builder and the commit; with beta's Evidence by the verifier it exits 0. TestBuildCheck covers the same rules, including an Actor: trailer by ID. All passed.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/build-check$'; ./internal/api -run TestBuildCheck (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-check](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-check.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
