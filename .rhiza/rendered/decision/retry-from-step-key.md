# retry-from-step-key

decision · `01M7QQ8DFNW9VMVDF26CX2B16H`

## reasoning

The intent says the workflow names the step; a per-step key keeps the default (the failed step) for every step that sets nothing, and one rule for every retry path keeps a manual retry and the daemon's the same

## statement

A workflow names the step a failed run restarts at with a step key retry_from naming another step of the same workflow (prefixed like on.* targets when the workflow is included by a workflow step); a step without it restarts at itself (after the existing walk back to the wait that feeds a failed check). Every retry honours it: rhiza run resume --retry, rhiza run start on a failed change set and the daemon's

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-retry-conduct-plan](../source/m40-retry-conduct-plan.md)

### In

- decided-by [workflow-retry-from-key](../capability/ship-through-rhiza/workflow-step/workflow-retry-from-key.md)
- decided-by [retry-restarts-at-named-step](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-restarts-at-named-step.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
