# plan-clean-round-critic-plan

source · `01M77NZ42MBQQ8QZ7SDDQEKYH0`

## quotes

```json
[]
```

## raw_sha256

96cbe7ec275501d729d33c271ba9e794acec8bb6b39f820c39e33c449305b06a

## session

dea9e255-b2d8-4d91-a3e9-01a1cf5cc817

## summary

Autonomous plan run for a bug request. Finding: in the built-in plan and plan-milestone workflows (internal/workflow/builtin/plan.toml, plan-milestone.toml), review always passes to critic, but rhiza run verdict (internal/api/run_review.go) refuses an empty verdicts list, so after a clean round (seen in run 01M77N1ENKEAY668AM2STW90V7, round 2) the critic can record nothing and cannot meet its prompt. The code review workflow's checks step already gives a clean round no critic (decision critics-on-serious-findings). One question was raised: skip the critic on a clean round, let run verdict accept an empty list, or both. With no owner, the arbiter decided plan-clean-round-skips-critic: route a clean round from review straight to fresh, and keep run verdict refusing an empty list. The change set serves m4-runs-hold-up, as the earlier fix review-visit-records-round did. One leaf was planned on plan-workflow.

## title

Plan of the fix: plan workflows send a clean review round to the critic

## Links

### In

- derived-from [plan-clean-round-skips-critic](../decision/plan-clean-round-skips-critic.md)
- derived-from [plan-clean-round-skips-critic-step](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-clean-round-skips-critic-step.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
