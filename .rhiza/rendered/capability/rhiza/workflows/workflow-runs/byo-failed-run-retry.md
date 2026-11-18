# byo-failed-run-retry

capability · `01M7RB62N25AX10VV3VWNZH77K` · parent [workflow-runs](../workflow-runs.md)

## acceptance

env -u'd of RHIZA_*/SHIP_*/REX_*: go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/byo-workflow-retry$' passes: the fake close-out fails on its first call, the run fails at close-out, the next daemon --once poll resumes it from close-out (not from services or work), and it reaches succeeded with intervention_count 0; (failure case) with the failed-retry budget spent (rhiza.daemon.watchdog.retryBudget or the failed-retry count set to 0) the run stays failed and rhiza inbox lists it

## files

cmd/rhiza/testdata/script/byo-workflow-retry.txtar, cmd/rhiza/main_test.go (processScripts entry with its reason); internal/api only for a gap found, fixed in the generic path

## must_not

Change docs/examples/byo/*.toml; special-case a workflow or step name; touch the network, Rex, a real agent, gh or docker; bind a unix socket

## statement

Under the daemon, a BYO example run whose close-out step fails once (a transient command failure) is retried by the watchdog from that failed step with no one acting, and ends succeeded with no journalled intervention

## Links

### Out

- depends-on [byo-workflow-daemon-e2e](byo-workflow-daemon-e2e.md)
- decided-by [byo-example-copy-synced](../../../../decision/byo-example-copy-synced.md)

### In

- verifies [byo-failed-run-retry-verified](../../../../evidence/byo-failed-run-retry-verified.md)
- verifies [byo-failed-run-retry-qa-1](../../../../evidence/byo-failed-run-retry-qa-1.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
