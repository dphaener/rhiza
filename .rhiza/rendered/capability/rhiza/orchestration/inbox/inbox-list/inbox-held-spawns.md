# inbox-held-spawns

capability · `01M7QGA6T7SB7NWYPE3NN8DYKS` · parent [inbox-list](../inbox-list.md)

## acceptance

go test -trimpath ./internal/api -run Inbox and a new testscript cmd/rhiza/testdata/script/inbox-spawn-held.txtar (go test -trimpath ./cmd/rhiza -run TestScripts/inbox-spawn-held): a live run whose journal has an open spawn-held entry gives a spawn-held item carrying its detail in rhiza inbox and rhiza inbox --json; (failure case) a spawn-held entry released by a later spawn or session entry, or in an ended run, gives no item; describe-golden regenerated (TestScripts/describe-golden)

## files

internal/api/inbox.go, internal/api/inbox_test.go, internal/api/web_needs_you.go, cmd/rhiza/testdata/script/inbox-spawn-held.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change the load-high, memory-low, disk-low or tmp-low machine-wide items; give one hold more than one item ID

## statement

rhiza inbox shows one spawn-held item for each spawn or step session a live run's journal holds, naming the change set, the step, the role of a spawned child, and the reading and limit from the entry's detail; the item is gone once the child or the session starts

## Links

### Out

- depends-on [step-session-waits-holds](../../../workflows/workflow-runs/step-session-waits-holds.md)
- decided-by [m40-spawn-gate-reads-daemon-state](../../../../../decision/m40-spawn-gate-reads-daemon-state.md)

### In

- verifies [inbox-held-spawns-verified](../../../../../evidence/inbox-held-spawns-verified.md)
- verifies [inbox-held-spawns-qa-1](../../../../../evidence/inbox-held-spawns-qa-1.md)
- verifies [inbox-held-spawns-qa-2](../../../../../evidence/inbox-held-spawns-qa-2.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
