# build-retries-verified

evidence · `01M5ZMT6S45GTS76DG7ZDX3AMQ` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

TestLeavesStepRetries: the verifier fails twice (findings, then exit) and then passes; it finds three builder sessions, distinct, on sonnet, opus, opus, each LeafSpec and RHIZA_ATTEMPT naming attempts 1, 2, 3, and the leaf passed at attempt 3. TestLeafPrompt finds the builder prompt rendered from that spec naming its attempt ('This is attempt 2 of 3 on leaf leaf-a'), and StartLeaf renders the prompt with it. TestLeavesStepRetriesCap: three verifier fails start no fourth builder and the leaf failed. TestLeavesStepRetriesBuilderFailCounts and TestLeavesStepRetriesResume also pass. Note: the prompt text is checked by the api unit test, not inside the engine test.

## tests

go test -count=1 ./internal/run -run TestLeavesStepRetries; ./internal/api -run TestLeafPrompt (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-retries](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-retries.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
