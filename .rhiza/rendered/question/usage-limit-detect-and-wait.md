# usage-limit-detect-and-wait

question · `01M6507JM2EC82K2P7MJHPK4Z9`

## statement

How does a run recognise the Claude usage limit, and how long does it wait before resuming: the reset time given in the limit message, or a fixed interval?

## Links

### Out

- blocks [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

### In

- answers [usage-limit-wait-from-message](../decision/usage-limit-wait-from-message.md)
- answers [usage-limit-detect-and-wait-answer](../decision/usage-limit-detect-and-wait-answer.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
