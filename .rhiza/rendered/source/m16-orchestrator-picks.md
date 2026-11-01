# m16-orchestrator-picks

source · `01M6BQJY7M0ZNYE2THSFGCZGNW`

## quotes

```json
[
  {
    "by": "owner",
    "text": "if you run out of work keep picking high value tasks. don't stop and wait for my input.",
    "where": "orchestrator session, 2026-10-31"
  }
]
```

## raw_sha256

994e035ce2a027fd110d52e6aeecfcee711e2eb6aacf2560995dea9bc156e897

## summary

With milestone 15 the last one queued, the orchestrator planned milestone 16 under the owner's standing instruction. Reviewing a change set today shows op types, IDs and handles only: change show prints one line per op, and the change-show app shows that output, so an owner (or arbiter) accepting a plan cannot see which statement changed or how. An update op stores only the changed fields (a delta) and the version it was written against, so the old value must come from the node at that base, not from the current as-is, which later landed work may have moved. There is no line-diff code in the repo. The founding decision field-level-diffs ('Diffs are structured and field-level') sits in the Area change-sets umbrella as change-diff; notes on a change need their own storage format and stay in the umbrella.

## title

Milestone 16: see what a change set changes

## Links

### In

- derived-from [m16-picked-by-orchestrator](../decision/m16-picked-by-orchestrator.md)
- derived-from [diff-against-op-base](../decision/diff-against-op-base.md)

## History

- `01M6BQJPJPWB8CZVCKAKWWC9T7` Plan milestone 16: see what a change set changes, picked by the orchestrator under the owner's standing instruction
