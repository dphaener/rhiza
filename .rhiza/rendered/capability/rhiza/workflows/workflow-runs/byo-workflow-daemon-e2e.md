# byo-workflow-daemon-e2e

capability · `01M7RB622RJPBYCSMFMQGDH35M` · parent [workflow-runs](../workflow-runs.md)

## acceptance

env -u'd of RHIZA_*/SHIP_*/REX_*: go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/byo-workflow-daemon$' passes: the daemon --once start names the myteam-flow run, the run reaches succeeded with each step in order, the fake docker logs 'compose down' from the cleanup, and rhiza run status --json shows intervention_count 0; (failure case) the first gate's fake make lint exits 1, the gate hands back to work with handback.md naming make lint, and the second gate passes; go test -trimpath ./cmd/rhiza -run 'TestDocsExample' passes and fails when a script's archived .rhiza/workflows/myteam-flow.toml differs from docs/examples/byo/workflow.toml

## files

cmd/rhiza/testdata/script/byo-workflow-daemon.txtar, cmd/rhiza/main_test.go (processScripts entry with its reason), cmd/rhiza/docs_examples_test.go

## must_not

Change docs/examples/byo/*.toml; special-case a workflow or step name in production code; touch the network, Rex, a real agent, gh or docker; bind a unix socket

## statement

A testscript runs the BYO example workflow myteam-flow under the daemon end to end: the daemon picks it by dispatch [ready] for an accepted change set and its services, work, gate, signoff, push, wait-ci, merge and close-out steps pass, the stop-services cleanup runs, and the run ends succeeded with no journalled intervention

## Links

### Out

- depends-on [open-pr-review-when-declared](../../../ship-through-rhiza/open-pr-review-when-declared.md)
- decided-by [byo-example-copy-synced](../../../../decision/byo-example-copy-synced.md)

### In

- depends-on [byo-spawn-hold-intervention](byo-spawn-hold-intervention.md)
- depends-on [byo-failed-run-retry](byo-failed-run-retry.md)
- blocks [byo-open-pr-needs-review](../../../../question/byo-open-pr-needs-review.md)
- verifies [byo-workflow-daemon-e2e-verified](../../../../evidence/byo-workflow-daemon-e2e-verified.md)
- verifies [byo-workflow-daemon-e2e-qa-1](../../../../evidence/byo-workflow-daemon-e2e-qa-1.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
