# printout-counts-used-tokens-verified

evidence · `01M65BR4E07E9V3YHEAAJ40AMB` · actor [claude-code-60ed4ab0](../actor/claude-code-60ed4ab0.md)

## result

PASS

## statement

At b706c7d ran go test ./cmd/rhiza -run TestScripts/printout-counts-used-tokens -count=1 -v: PASS. The script runs a three-step workflow with a fake claude whose transcript shows usage 10/20/30/40 plus helper 1/2/3/4, strips every tokens entry from the journal, then run status --json shows 3 tokens objects and totals session 30/60/90/120, helpers 3/6/9/12; human run status lists each session with its counts and the Totals line, with no 'no token use recorded'; the journal still has no tokens entry afterwards and fold-ok passes. Diff read: runInfo falls back to transcriptTokens (session record -> adapter -> readTranscript -> TokenUse) when the journal has no reading and adds it to totals; the final printout (runEngineRetry's RunResponse) and run status, resume, abandon and signoff all build through that same runInfo, so the final printout gets the same counts. go test ./... and fmt/vet pass; make self-check fails only on rendered files for actors outside this commit (claude-code-60ed4ab0, this verifier, and claude-code-afe0a6eb).

## Links

### Out

- verifies [printout-counts-used-tokens](../capability/rhiza/runtime/session-journal/run-status-tokens/printout-counts-used-tokens.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
