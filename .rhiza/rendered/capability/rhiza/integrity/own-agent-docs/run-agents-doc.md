# run-agents-doc

capability · `01M76P0B7W8M843KJ8BFHR6C80` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

A test in cmd/rhiza reads agents/run-agents.md and finds each lesson (trailers in one paragraph, already-landed change set, empty commit, merge origin/main, ship json after a hand push); it and the agent-docs consistency test pass with go test -trimpath ./cmd/rhiza -run on their names

## statement

agents/run-agents.md, listed in agents/README.md, holds the lessons for builders, verifiers, checks and fixers working on Rhiza itself: Leaf:, Change: and Actor: trailers share the last paragraph with Co-Authored-By, run steps writing to the model handle a run change set that already landed, findings HEAD already fixes still get a new (empty) commit, a conflicting PR merges origin/main and is never rebased and force-pushed, and a hand push to a live run PR updates .git/rhiza/ship/<cs>.json

## Links

### In

- verifies [run-agents-doc-verified](../../../../evidence/run-agents-doc-verified.md)
- verifies [run-agents-doc-qa-1](../../../../evidence/run-agents-doc-qa-1.md)
- verifies [pr-497-role-docs-for-rhiza-agents](../../../../evidence/pr-497-role-docs-for-rhiza-agents.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
