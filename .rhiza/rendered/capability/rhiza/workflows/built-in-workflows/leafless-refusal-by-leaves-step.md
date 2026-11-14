# leafless-refusal-by-leaves-step

capability · `01M7AHXFBEANQM6G4MHTGZJ9MA` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/api -run TestCheckShipLeaves shows a leafless change set refused for ship-change-classic and for a repo workflow that includes build, and accepted for a workflow with no leaves step; go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-needs-leaves passes with the script pinned to ship-change-classic. Touches internal/api/run.go, internal/api/run_test.go, cmd/rhiza/testdata/script/run-ship-needs-leaves.txtar.

## statement

rhiza run start refuses a change set with no leaves for any workflow whose expanded steps include a leaves step (ship-change-classic, or a repo workflow including build), instead of only for a workflow named ship-change, so a workflow without one, like the conductor default, starts on a leafless change set.

## Links

### Out

- decided-by [conductor-plans-leaves](../../../../decision/conductor-plans-leaves.md)

### In

- depends-on [classic-tests-pinned](classic-tests-pinned.md)
- verifies [leafless-refusal-by-leaves-step-verified](../../../../evidence/leafless-refusal-by-leaves-step-verified.md)
- verifies [leafless-refusal-by-leaves-step-qa-2](../../../../evidence/leafless-refusal-by-leaves-step-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
