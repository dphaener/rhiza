# verifier-reruns-reported-checks-verified

evidence · `01M788EECFT3WP9E9EWNBJ74WW` · actor [claude-code-66b6cdcf](../actor/claude-code-66b6cdcf.md)

## result

PASS

## statement

At a62c7cd14: `go test -trimpath -count=1 ./internal/guide -run Build -v` exits 0, TestBuildVerifierRerunsChecks passes, asserting build.md's Verifier section says to re-run every command in the commit's `Checks:` block and every acceptance command, and that a different exit code or no `Checks:` block is a finding naming the command (read in the diff of 0ce75a449). Re-ran a62c7cd14's reported Checks: `make self-check` in a clean detached worktree with RHIZA_ROOT/CHANGE/RUN unset exits 0 (render no longer differs; attempt 1's finding f1 is fixed); the guide test also exits 0 as reported.

## Links

### Out

- verifies [verifier-reruns-reported-checks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/verifier-reruns-reported-checks.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
