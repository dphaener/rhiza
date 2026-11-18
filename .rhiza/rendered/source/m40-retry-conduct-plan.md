# m40-retry-conduct-plan

source · `01M7QQ77593WP7HYTQ9VQRFHYW`

## quotes

```json
[
  {
    "by": "orchestrator",
    "text": "plan-bug run 01M7PSG0M7 passed review on round 2 and merged, then failed at plan/close-out on a transient `git pull --ff-only origin main` error in the main checkout ... It stayed failed until the orchestrator's `rhiza run start` on the change set retried close-out, which passed. For m40-retry: a transient command failure in a cleanup step is the case an automatic retry of the failed step should cover with no one watching.",
    "where": "hands-off-milestone-interview.md, later evidence 2026-11-18"
  },
  {
    "by": "orchestrator",
    "text": "today rhiza run start on a failed change set already retries the failed run at its failed step",
    "where": "machine rules message to the conductor of run 01M7QQ2WVR"
  }
]
```

## raw_sha256

e7222e33f280cca01e8f37952f64ff727f1149bb44377cb07af023a5e90ff725

## summary

The conductor of run 01M7QQ2WVR planned change set 01M7NVMM (m40-retry) with no owner present, from the m40-retry intent and its Decision m40-retry-right-step (the owner's pick 'Right step + findings', relayed by the orchestrator), and the orchestrator's later evidence in hands-off-milestone-interview.md. Research: run.Engine.Retry reopens a failed run at RetryStep(def, st, failed), which only walks back to a wait feeding the failed check; no workflow key names a restart step. The watchdog's watchdogTransient retries a failed autonomous run only when the failure classifier calls it transient (budget rhiza.daemon.watchdog.retryBudget) and leaves every real or unclassified failure to the inbox. A step session's brief.md (run.RenderBrief via stepRunner.fillSession) carries nothing of an earlier failed attempt; handback.md is only written on a hand-back. Built-in plan.toml's review loop (review, has-findings, critic, fresh, check, revise, rulings, apply-rulings) ends the run at its own step when a cap is reached, so a retry today re-enters that step instead of planning again. No owner words: autonomous run.

## title

Retry change set plan, autonomous conduct run 01M7QQ2WVR

## Links

### In

- derived-from [retry-from-step-key](../decision/retry-from-step-key.md)
- derived-from [plan-review-retries-at-interview](../decision/plan-review-retries-at-interview.md)
- derived-from [retry-brief-previous-attempt](../decision/retry-brief-previous-attempt.md)
- derived-from [daemon-failed-retries-setting](../decision/daemon-failed-retries-setting.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
