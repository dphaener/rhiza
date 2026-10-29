# inbox-stale-interview

source · `01M62S4SD2FSQ9FHXDB7SY9YMQ`

## quotes

```json
[
  {
    "by": "agent (accepted by the owner)",
    "text": "Failed and gave-up runs (Recommended): both leave the inbox once their change set is landed, withdrawn or rejected. Questions and sign-offs stay listed.",
    "where": "interview: stale items"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Yes, stale (Recommended): a failed or gave-up run on a superseded change set leaves the inbox, like one on a landed, withdrawn or rejected change set.",
    "where": "interview: superseded"
  }
]
```

## raw_sha256

0ca6b95fb5f86969fdf79a6a0c37d1eb40fbf3c1734ce56aad482327d88aebb8

## session

29b4ed69-1d8a-481f-a616-01a185858ae8

## summary

Planned with the owner on change set 01M62MKZFRPZ2C7QTYAF815XAG. A change set is stale when it has landed, been withdrawn, rejected or superseded. The owner picked: failed runs and runs the daemon gave up on leave the inbox once their change set is stale, questions and sign-offs stay listed, and superseded counts as stale. Still to build: rhiza run abandon takes a failed run, and the inbox names no next command that can only fail again. Nothing is open.

## title

Inbox clears stale failures, with the owner

## Links

### In

- derived-from [answer-29b4ed69-1](../decision/answer-29b4ed69-1.md)
- derived-from [answer-29b4ed69-2](../decision/answer-29b4ed69-2.md)
- derived-from [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)
- derived-from [abandon-failed-run](../capability/rhiza/workflows/workflow-runs/abandon-failed-run.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
