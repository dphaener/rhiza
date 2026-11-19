# land-commit-qa-1

evidence · `01M7RM16DA8Z7P1FMJA78TEKSQ` · actor [claude-code-5f8b9e60](../actor/claude-code-5f8b9e60.md)

## result

PASS

## statement

QA: built binary in scratch git repos: land --commit commits only .rhiza with Change:/Actor: trailers, staged code.go stays staged and uncommitted; draft exits 3 and dry-run exits 2 with HEAD unchanged; RHIZA_ACTOR unset falls back to the user.email actor; plain land unchanged; via a fake head build --commit is not forwarded and the commit still carries trailers; land-commit and describe-golden scripts and TestLandCommitViaHead pass. Minor: outside a git repo it lands, then fails at git add

## Links

### Out

- verifies [land-commit](../capability/rhiza/change-sets/change-lifecycle/land-commit.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
