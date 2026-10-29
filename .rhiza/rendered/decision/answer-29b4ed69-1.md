# answer-29b4ed69-1

decision · `01M62S4WQPRZYZFYRW27JXDD2N`

## options

```json
[
  "Failed and gave-up runs (Recommended)",
  "Failed runs only",
  "Every kind of item"
]
```

## picked

Failed and gave-up runs (Recommended)

## question

Which inbox items leave the list once their change set has landed, been withdrawn or rejected? The intent names failed runs only. A gave-up run's next command, rhiza run resume, is a dead end on such a change set too.

## session

29b4ed69-1d8a-481f-a616-01a185858ae8

## statement

Which inbox items leave the list once their change set has landed, been withdrawn or rejected? The intent names failed runs only. A gave-up run's next command, rhiza run resume, is a dead end on such a change set too. — Failed and gave-up runs (Recommended)

## wording

the agent's option, picked by the owner in session 29b4ed69-1d8a-481f-a616-01a185858ae8

## Links

### Out

- derived-from [inbox-stale-interview](../source/inbox-stale-interview.md)

### In

- decided-by [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
