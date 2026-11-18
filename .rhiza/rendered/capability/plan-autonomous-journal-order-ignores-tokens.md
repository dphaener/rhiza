# plan-autonomous-journal-order-ignores-tokens

capability · `01M7RFQFXB7SZSEMAB7X7X4M0W`

## acceptance

Every go test below runs inside systemd-run --user --scope --quiet -p MemoryMax=6G -p TasksMax=300 with -p 2 and RHIZA_*/SHIP_*/REX_* unset, one at a time. (1) go test -trimpath -p 2 -count=3 ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous$' passes. (2) (failure case) With two temporary scratch edits, never committed: time.Sleep(1500 * time.Millisecond) right after r.log("step %s: %s", step.Name, res.outcome) in internal/run/engine.go's agent-step end, and 'sleep 0.8' between the fake planner's 'rhiza session done' and its 'exit 0' in the script's bin/claude, the script at main f7313222e fails at the journal-order line with 'unblocked tokens tokens ended tokens entered' in the journal (as in PR #655's CI), and with the fix it passes on 3 runs; both edits reverted and the commit's Checks: block records both outcomes. (3) (failure case) The new tokens assertion fails when its step name is changed in scratch to one with no tokens entry (e.g. step nosuch), showing it is not vacuous; reverted. (4) Every other assertion of the script is unchanged, and go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous' (the -left and -unapplied scripts too) passes.

## files

cmd/rhiza/testdata/script/plan-workflow-autonomous.txtar

## must_not

Change any Go code (internal/ or cmd/rhiza/*.go), including token journaling order; quarantine, skip or add the script to any skip list; loosen any assertion other than the tokens placement in the order check (the rest of the order, from started to the first entered after ended, stays exact); edit other scripts; use busy loops, CPU stress or parallel stress runs to reproduce

## statement

In cmd/rhiza/testdata/script/plan-workflow-autonomous.txtar the journal-order check (the grep at line 66 and its stdout regex at line 68) leaves tokens entries out of the sequence it matches, so it reads 'started process entered session blocked asked arbiter arbitrated answered unblocked ended entered ' whatever position the tokens entries take, and a separate assertion requires at least one tokens entry with step interview in the run journal

## Links

### Out

- decided-by [plan-autonomous-tokens-order-free](../decision/plan-autonomous-tokens-order-free.md)

### In

- verifies [plan-autonomous-journal-order-ignores-tokens-verified](../evidence/plan-autonomous-journal-order-ignores-tokens-verified.md)
- verifies [plan-autonomous-journal-order-ignores-tokens-qa-1](../evidence/plan-autonomous-journal-order-ignores-tokens-qa-1.md)

## History

- `01M7RFDAZ53YZ3775D7SV70Z4G` Fix flaky test TestScripts/plan-workflow-autonomous: failed then passed on a rerun in PR #655 at head 1be14770b82366ab1535a1b38a36a890e8b75939; reproduce it locally, repeatedly, then fix it or quarantine it
