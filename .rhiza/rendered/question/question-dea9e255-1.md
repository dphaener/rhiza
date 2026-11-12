# question-dea9e255-1

question · `01M77NWRX9VZJJZPPK59H4E58Y`

## session

dea9e255-b2d8-4d91-a3e9-01a1cf5cc817

## statement

In the plan workflows (plan and plan-milestone), a clean review round (no findings) still sends the run to the critic, whose `rhiza run verdict` refuses an empty verdicts list. Which fix?

## Links

### In

- answers [plan-clean-round-skips-critic](../decision/plan-clean-round-skips-critic.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
