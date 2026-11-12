# release-marks-verified

evidence · `01M78298RWM36VXSEC1Z8WBN7H` · actor [claude-code-3dd70c3b](../actor/claude-code-3dd70c3b.md)

## result

PASS

## statement

At 605393a64: go test -trimpath -count=1 ./tools/release/... passes (TestReadMarks fixture log: breaking: and case-insensitive Severe: notes discussing change sets are read with their reasons; plain notes, 'breaking changes' without colon, notes on non-change-set nodes, retired notes and notes with retired links are ignored). End to end, in a scratch repo built from this commit, three real `rhiza note add <cs>` notes (breaking:, severe:, plain) gave `tools/release marks` = {cs:{breaking:[drops --old],severe:[loses data]}}. gofmt clean, go vet ok. make self-check passes in a clean detached worktree of the commit (RHIZA_* env cleared; in the run worktree it only failed on the verifier's own untracked actor render). Guidance: CLAUDE.md gains a Release marks section telling planners (plan-bug included) when to write breaking:/severe: and reviewers to add missed marks; plan-bug.toml and bugfix.toml carry the same guidance as comments.

## Links

### Out

- verifies [release-marks](../capability/rhiza/distribution/release-marks.md)

## History

- `01M76B077D23CJ28C690Y166VY` Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against
