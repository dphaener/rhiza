# answer-review-01m6fdsp-r3-f1

decision · `01M6FGJ2RFDF73NV08DR16EB3W`

## options

```json
[
  "Run plans headless (Recommended)",
  "Keep plan tabs"
]
```

## question

Plan review round 3 left finding f1: request plan runs still start in a Rex tab. How should it be settled?

## quote

Run plans headless (Recommended): widen daemon-headless-run-drivers so a request's plan run also starts in the background, and drop 'plan-request and run tabs' from daemon-repo-session. Only the plan's agent sessions open tabs.

## quote_where

the orchestrator session, 2026-11-03, owner's answer by choice

## Links

### Out

- answers [review-01m6fdsp-r3-f1](../question/review-01m6fdsp-r3-f1.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
