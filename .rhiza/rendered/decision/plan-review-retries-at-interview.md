# plan-review-retries-at-interview

decision · `01M7QQ8DHFZM683VRPQ000RVAD`

## reasoning

The intent's 'a plan whose review-check refuses restarts at plan': the run fails at whichever review-loop step reaches its cap, so each of them must name the interview

## statement

The built-in plan and plan-milestone workflows set retry_from = "interview" on every step of their review loop (review, has-findings, critic, fresh, check, revise, rulings, apply-rulings), so a plan whose review-check refuses is retried from its interview, which is the plan step of a workflow that includes it (plan-bug's plan/interview)

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-retry-conduct-plan](../source/m40-retry-conduct-plan.md)

### In

- decided-by [retry-restarts-at-named-step](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-restarts-at-named-step.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
