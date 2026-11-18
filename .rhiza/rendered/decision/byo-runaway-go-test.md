# byo-runaway-go-test

decision · `01M7RB620T9FT6Y81X278EJEWX`

## statement

Stall and runaway detection on a BYO workflow's steps are shown by Go tests in internal/api over a BYO-shaped workflow (agent step work, gate step gate), not a testscript: the process table reader is nil under RHIZA_TEST_NOW

## wording

(planner, from the code: watchdogProcTable in internal/api/watchdog_runaway.go returns nil in tests)

## Links

### Out

- derived-from [byo-daemon-plan-research](../source/byo-daemon-plan-research.md)

### In

- decided-by [byo-watchdog-stall-runaway](../capability/rhiza/workflows/workflow-runs/byo-watchdog-stall-runaway.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
