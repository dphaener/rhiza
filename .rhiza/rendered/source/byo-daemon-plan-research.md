# byo-daemon-plan-research

source · `01M7RB4B3DK2YBCWWGEKRE7ZGV`

## quotes

```json
[
  {
    "by": "orchestrator, summarising the owner",
    "text": "Whatever the daemon enforces (limits, stall recovery, retry, hygiene) must apply to any workflow's steps (docs/bring-your-own-agent.md, docs/examples/byo/), not only ship-change, plan-bug and bugfix. A done point should show a BYO workflow run end to end under the daemon.",
    "where": "~/orchestrator/hands-off-milestone-interview.md, Bring-your-own workflows"
  }
]
```

## raw_sha256

c274a9398f942544b73a359b2f9535606c36f6765c4c54dee871507af6de4e17

## session

2debf9c0-0390-4456-af6c-01a1f0aeb98a

## summary

Autonomous plan run (no owner). Research of the code for change set 01M7NVN: no test runs a bring-your-own workflow end to end under the daemon. Closest templates: cmd/rhiza/testdata/script/bugfix-workflow-autonomous.txtar (daemon --once with RHIZA_TEST_DETACH_LOG, dispatch [ready], fake rex/claude) and run-ship-change.txtar (fake gh for open-pr, ship check, merge); run-cleanup-steps.txtar for [[cleanup]]; workflow-invoke.txtar for installing a team skill under $HOME/.claude/skills. Spawn and step-start holds (run_spawn.go, engine.go StartHold), failed-run retry (daemon_watchdog.go watchdogTransient, FailedStep) and the intervention journal (internal/run/journal.go TypeIntervention) are not tied to built-in step names. Stall limits switch on step kind (watchdog_policy.go StallLimit): a gate step with no stall= gets the agent retry limit, not the command limit; the built-in gates set stall=70m, the BYO example workflow.toml gate sets none. Runaway detection (watchdog_runaway.go) cannot be triggered in a testscript (watchdogProcTable is nil under RHIZA_TEST_NOW), so it is covered by Go tests. cmd/rhiza/docs_examples_test.go already runs workflow check on docs/examples/byo/*.toml; a testscript cannot read repo files, so an archived copy needs a Go test holding it byte-identical to docs/examples/byo/workflow.toml.

## title

BYO workflow under the daemon: planning research

## Links

### In

- derived-from [byo-example-copy-synced](../decision/byo-example-copy-synced.md)
- derived-from [byo-gate-command-limit](../decision/byo-gate-command-limit.md)
- derived-from [byo-runaway-go-test](../decision/byo-runaway-go-test.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
