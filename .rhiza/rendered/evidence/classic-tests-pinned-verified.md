# classic-tests-pinned-verified

evidence · `01M7C374RWSX36CWMR1H9NYGKA` · actor [claude-code-96f0b235](../actor/claude-code-96f0b235.md)

## result

PASS

## statement

At 29632c774: go test -trimpath ./cmd/rhiza -run 'TestScripts/(<the 30 named scripts>)$' passed all 30 (exit 0); go test -trimpath ./internal/api -run 'TestPrePushJournal|TestInboxReadsChangeSetsOnce' passed; go test -trimpath ./internal/workflow -run with the eight ship tests passed (all eight PASS). grep -L ship-change-classic on the 30 scripts, ship_prepush_test.go, inbox_fold_test.go and workflow_test.go printed nothing; grep -nE 'Builtin\("ship-change"\)|load\("ship-change"\)' internal/workflow/workflow_test.go printed nothing. The diff shows run-start scripts pass --workflow ship-change-classic, daemon-headless-runs copies a dispatch.toml with [ready] feature = "ship-change-classic", seeded started entries and run.Started name ship-change-classic, and the workflow_test.go hunks are only inside the eight named tests. None of the 12 shadowing scripts was touched, and nothing outside the leaf changed except the builder's staged actor record. No TODOs. The Checks rhiza ship generate and rhiza render both exited 0 with no tracked change.

## Links

### Out

- verifies [classic-tests-pinned](../capability/rhiza/workflows/built-in-workflows/classic-tests-pinned.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
