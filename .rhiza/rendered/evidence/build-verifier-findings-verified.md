# build-verifier-findings-verified

evidence · `01M5ZMT6F19NTKP4RRYK58SA2A` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Testscript build-verifier-findings: a verifier runs rhiza run findings --leaf twice and the journal holds findings f1 and f2 against that leaf's attempt 1; rhiza run status shows them per attempt in text and --json and opens no review round; refusals (outsider leaf, no attempt yet, unknown field, empty reason, no findings) exit 2. Engine clause: TestLeavesStepVerifierEvidence finds the attempt failed for a verifier ending with neither Evidence nor findings, and for Evidence written by the builder's actor; the testscript shows the real check naming the builder's actor for leaf second. TestLeavesStepVerifierFindings finds journalled findings fail the attempt. All passed.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/build-verifier-findings$'; ./internal/run -run 'TestLeavesStepVerifierEvidence|TestLeavesStepVerifierFindings' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-verifier-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-verifier-findings.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
