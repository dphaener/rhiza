# model-lints-plan-interview

source · `01M7AZR0EQN8T0NZSV96NBC055`

## quotes

```json
[]
```

## raw_sha256

2e299bcb1e94b114ac2c508069152249d647cb0ea064e493e5a9a4fb97993e86

## session

c0130112-6d9f-43e5-ab34-01a1d5f666af

## summary

Autonomous plan run for change set 01M76XTF (Model lints in review), with no owner present. Research: perf report rec 4 (model check fails cycle 1 in 37-42% of runs, mostly mechanical drift); a sample of model-check majors from the run journals showed acceptance naming removed tests or files, Decisions recorded but never applied, and capabilities whose code was deleted but not retired. rhiza run code-check (internal/api/run_codecheck.go) runs rhiza check and build-check, and the builtin review's code step has no on.fail, so a code-check failure ends the run. The model check's prompt is built in internal/api/run_checks.go (StartCheck), and the check-model topic forbids reading the run's status. The planner asked three questions: where a lint hit goes, which nodes the stale-reference lint covers, and when a Decision counts as never applied. The arbiter answered them all with Decision lint-hits-feed-model-check, worded (assumed, would have escalated): lint hits never fail code-check; they are journalled on the code-check entry and handed to the cycle-1 model check as leads, which it confirms or drops; the stale-reference lint covers only the change set's own slice; a Decision is unapplied when the change set adds or changes it and no live node links decided-by to it (an answers-only Decision is exempt), or every node that does is a leaf of this change set with no Leaf: commit on the branch. The plan has four leaves: the stale-reference lint, the unapplied-Decision lint, code-check running and journalling them, and the model check receiving them as leads.

## title

Model lints in review, planned autonomously

## Links

### In

- derived-from [lint-hits-feed-model-check](../decision/lint-hits-feed-model-check.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
