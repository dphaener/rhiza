# m20-orchestrator-picks

source · `01M6C8ED7SZSQG897MA7CB9D2B`

## quotes

```json
[
  {
    "by": "owner",
    "text": "just checking in. keep cracking on",
    "where": "orchestrator session, 2026-11-01"
  }
]
```

## raw_sha256

c3c33b41ddbf5ce40dab9d7987bc8416b6a4288ec6cc21d5f0f9936eafa25f7d

## summary

With the queue running one change set at a time (each regenerates describe-golden), the orchestrator picked work for the idle daemon slot that changes no request type: the older proposed change set 01M5SWTFB3, what-if reads under the 50 ms read budget. Milestone 15's what-if change cut the 50k-node what-if show from about 1,620 ms to about 240 ms by folding only the in-flight change sets over the stored as-is rows, but it still builds a full in-memory model of the whole as-is (about 110 ms on the profile before, the rest now), and TestReadBudget times only as-is and to-be reads. Meeting the budget needs the reader to answer from the index's stored as-is rows with an overlay of only the records the change set and its closure touch or create, rather than a full model.

## title

Milestone 20: what-if reads meet the budget

## Links

### In

- derived-from [m20-picked-by-orchestrator](../decision/m20-picked-by-orchestrator.md)
- derived-from [whatif-overlay-reader](../decision/whatif-overlay-reader.md)

## History

- `01M6C8EBWCSJZC8TGSN36SZ3QR` Plan milestone 20: what-if reads meet the budget, picked by the orchestrator under the owner's standing instruction
