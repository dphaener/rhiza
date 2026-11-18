# byo-spawn-hold-intervention

capability · `01M7RB62EDVAPFFYGG892NFT57` · parent [workflow-runs](../workflow-runs.md)

## acceptance

env -u'd of RHIZA_*/SHIP_*/REX_*: go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/byo-workflow-hold-intervention$' passes: with RHIZA_TEST_LOAD1/5 over maxLoad the work step does not start and rhiza inbox shows a held step session for the myteam-flow run; once the load drops the work session starts; a rhiza session message from outside the run to the work session leaves rhiza run status --json with intervention_count 1 naming the message; (failure case) a message the run's own session sends records no intervention

## files

cmd/rhiza/testdata/script/byo-workflow-hold-intervention.txtar, cmd/rhiza/main_test.go (processScripts entry if it starts a process); internal/api only for a gap found, fixed in the generic path

## must_not

Change docs/examples/byo/*.toml; special-case a workflow or step name; touch the network, Rex, a real agent, gh or docker; bind a unix socket

## statement

Under the daemon, the BYO example's work step session is held by the load hold like a built-in step's (held while load is high, started once it drops) and an outside rhiza session message to that session is journalled as an intervention shown by run status

## Links

### Out

- depends-on [byo-workflow-daemon-e2e](byo-workflow-daemon-e2e.md)
- decided-by [byo-example-copy-synced](../../../../decision/byo-example-copy-synced.md)

### In

- verifies [byo-spawn-hold-intervention-verified](../../../../evidence/byo-spawn-hold-intervention-verified.md)
- verifies [byo-spawn-hold-intervention-qa-1](../../../../evidence/byo-spawn-hold-intervention-qa-1.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
