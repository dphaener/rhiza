# build-workflow-file-verified

evidence · `01M5ZMT8BN7SD8X50QZ83CVB0M` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Testscript build-workflow runs rhiza run start --workflow build (built-in, no workflow file) with fake rex and claude on an accepted change set of two leaves, second depending on first: first is built and verified before second's builder, build-check passes, the run is succeeded, and the change set is still accepted. describe-golden passes against describe --json. TestBuildBuiltin checks the built-in's steps, models and no land. All passed.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/(build-workflow|describe-golden)$'; ./internal/workflow -run TestBuildBuiltin (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-workflow-file](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-workflow-file.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
