# plan-clean-round-skips-critic

decision · `01M77NY73VZ2J0F7N50ZD6SSXS`

## reasoning

The change set's intent offers two fixes: route a clean round from review to fresh, or let run verdict accept an empty list. Routing is the better fit. The code review workflow (internal/workflow/builtin/review.toml) already gives a critic only to a check with a blocker or major, so a clean round skipping the critic matches existing behaviour. A critic has nothing to judge in an empty round, so starting an agent session for it only spends a session, and an empty verdict list would be a record that says nothing. run verdict's rule that it needs at least one verdict, each naming a finding of the latest round (internal/api/run_review.go), stays a useful guard against a critic that records nothing for a round that has findings. Loosening it would also weaken that guard for the code review critics. plan/fresh and plan/check then judge the clean round as before (review-check passes a round with no confirmed findings), so only the plan workflows' routing changes: a command step or condition between review and critic that sends clean rounds to fresh.

## statement

In the plan and plan-milestone workflows, a review round with no findings skips the critic: after plan/review, the run goes to the critic only when the latest round has findings, and a clean round goes straight to plan/fresh. `rhiza run verdict` keeps refusing an empty verdicts list.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-dea9e255-1](../question/question-dea9e255-1.md)
- derived-from [plan-clean-round-critic-plan](../source/plan-clean-round-critic-plan.md)

### In

- decided-by [plan-clean-round-skips-critic-step](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-clean-round-skips-critic-step.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
