# default-tests-follow-conductor-verified

evidence · `01M7C62HK4J232BPXNRSNE06V2` · actor [claude-code-980d7c9c](../actor/claude-code-980d7c9c.md)

## result

PASS

## statement

At 15987daa1 (attempt commits ecc7874d4 and 15987daa1): go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-start-ready-dispatch|change-workflow-pin|bugfix-workflow-autonomous|daemon-root-ship-workflow)$' ok; go test -trimpath ./internal/daemon -run 'TestResumable|TestPick|TestSkips' ok. Diff read: run-start-ready-dispatch and change-workflow-pin now expect the default's dry run on a leafless change set to succeed with '- Workflow: ship-change' (old leaves-step refusals removed); pick_test.go moves the land-step resumable cases to ship-change-classic and adds ship-change cases on the new builtin's steps (merge/main-ci/close-out resumable, conduct/gate/push not), matching workflow describe ship-change of this build; bugfix-workflow-autonomous and daemon-root-ship-workflow use repo workflows and need no change. No test deleted, no TODOs, only the leaf's files plus .rhiza records changed.

## Links

### Out

- verifies [default-tests-follow-conductor](../capability/rhiza/workflows/built-in-workflows/default-tests-follow-conductor.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
