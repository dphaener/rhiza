# default-tests-follow-conductor

capability · `01M7AMSA5N0PTMBVTD00XHHBSG` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-start-ready-dispatch|change-workflow-pin|bugfix-workflow-autonomous|daemon-root-ship-workflow)$' passes with ship-change the conductor default, run-start-ready-dispatch and change-workflow-pin showing the default's dry run on a leafless change set succeeding with '- Workflow: ship-change'; go test -trimpath ./internal/daemon -run 'TestResumable|TestPick|TestSkips' passes. Touches cmd/rhiza/testdata/script/run-start-ready-dispatch.txtar, cmd/rhiza/testdata/script/change-workflow-pin.txtar, cmd/rhiza/testdata/script/bugfix-workflow-autonomous.txtar, cmd/rhiza/testdata/script/daemon-root-ship-workflow.txtar and internal/daemon/pick_test.go.

## statement

The existing tests of the default ship-change itself follow the conductor default: run-start-ready-dispatch.txtar and change-workflow-pin.txtar expect the default's dry run to start a leafless change set (no 'ship-change builds and reviews its leaves' or leaves-step refusal); bugfix-workflow-autonomous.txtar and daemon-root-ship-workflow.txtar assert the new builtin's steps; and the daemon pick tests in internal/daemon/pick_test.go that assert a step of the default name the new builtin's steps. No test is deleted.

## Links

### Out

- depends-on [builtin-conductor-default](builtin-conductor-default.md)
- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)

### In

- verifies [default-tests-follow-conductor-verified](../../../../evidence/default-tests-follow-conductor-verified.md)
- verifies [default-tests-follow-conductor-qa-2](../../../../evidence/default-tests-follow-conductor-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
