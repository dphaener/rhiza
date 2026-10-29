# answer-29b4ed69-2

decision · `01M62S4X7A6ZP1YT7V399VCNTF`

## options

```json
[
  "Yes, stale (Recommended)",
  "No, stays listed"
]
```

## picked

Yes, stale (Recommended)

## question

Does a superseded change set also count as stale, alongside landed, withdrawn and rejected? Superseded is a terminal state the intent does not name, and a run on it can only be retried into a change set that has moved on.

## session

29b4ed69-1d8a-481f-a616-01a185858ae8

## statement

Does a superseded change set also count as stale, alongside landed, withdrawn and rejected? Superseded is a terminal state the intent does not name, and a run on it can only be retried into a change set that has moved on. — Yes, stale (Recommended)

## wording

the agent's option, picked by the owner in session 29b4ed69-1d8a-481f-a616-01a185858ae8

## Links

### Out

- derived-from [inbox-stale-interview](../source/inbox-stale-interview.md)

### In

- decided-by [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
