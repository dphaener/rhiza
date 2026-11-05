# timing-budgets-already-done

source · `01M6J8C2YT50CPNYCXZYHMH1C0`

## quotes

```json
[]
```

## raw_sha256

2b9821da48bc7b624e70b140a37b53858f842b31fc224773012123446d1a5741

## session

89b23951-e5f9-4ffb-a6b5-01a1a484dffc

## summary

Autonomous plan run. The orchestrator reported the request already done: PR #225 (change set 01M6HR3VBT4D, landed) makes make test skip TestReadBudget|TestBudget|TestGeneratedModel in go test ./... and run them afterwards alone with -p 1 (capability budgets-measure-in-quiet, Decision budgets-run-alone). The Makefile on this branch confirms it. The orchestrator ruled: plan no new code; plan this change set as a no-op related to 01M6HR3VBT4D. No leaves were written.

## title

Plan of Timing budgets measure in: already done

## History

- `01M6J89D79MC3GB2P5EKFNDZ3G` Timing budgets measure in quiet: TestReadBudget, TestBudget (internal/slice) and modelapp's TestGeneratedModel time reads against fixed ~50ms limits while rhiza.test's scripts still run beside them (quiettest.Wait gives up after 3 minutes, cmd/rhiza takes 5-11), so CI fails them on noisy runners (PR #210 failed twice on the same SHA that passed a third time) and a loaded VM fails them every time. Run the budget tests on their own after go test ./... in make check (or make quiettest wait for rhiza.test to finish), so the budgets measure reads, not neighbours
