# pr-486-fix-plan-workflows-visit

evidence · `01M77X8D7N8X5VPE1KDQ1F7YWH` · actor [dphaener](../actor/dphaener.md)

## head

6f30e305c976ad85408b03141575219424557150

## pr

https://github.com/dphaener/rhiza/pull/486

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #486 (https://github.com/dphaener/rhiza/pull/486) passed CI at 6f30e30

## status

verified

## Links

### Out

- verifies [plan-clean-round-skips-critic-step](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-clean-round-skips-critic-step.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
