# conductor-land-commit-trailers

decision · `01M7RJTNS07WBZE13YKPYFF264`

## statement

The conductor's land commit is recognised by trailers rhiza writes, not by a new attribution rule: rhiza land --commit lands, renders and commits .rhiza with Change: and Actor: trailers (the actor every write resolves, the conductor session's own), and the conduct guide makes it the conductor's land step; handCommits' rule (Decision hand-commit-by-trailer) is unchanged, so a person's commit with only a Change: trailer still counts as a hand commit

## Links

### Out

- derived-from [false-hand-commit-plan-research](../source/false-hand-commit-plan-research.md)

### In

- decided-by [land-commit](../capability/rhiza/change-sets/change-lifecycle/land-commit.md)
- decided-by [conductor-land-commit-not-hand-commit](../capability/rhiza/workflows/workflow-runs/conductor-land-commit-not-hand-commit.md)
- decided-by [conduct-guide-land-commit](../capability/rhiza/guide/agent-guide/conduct-guide-land-commit.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
