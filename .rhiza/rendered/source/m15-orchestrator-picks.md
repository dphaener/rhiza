# m15-orchestrator-picks

source · `01M6BH834ZR9NRBBZWMJD3J96K`

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

ea14446f3ee6483a1cf59481d32a0742de8ec5818ee33a1c77ca4b316a279ba2

## summary

With milestone 14 the last one queued, the orchestrator planned milestone 15 under the owner's standing instruction. The new feed is drowned by actor records: of its first 1,000 items, 869 are 'Record agent actor' record change sets that write a single actor node, so filtering on kind alone would be wrong (records also carry plans and links). And every what-if read folds the whole log: Index.WhatIf runs foldStored over all entries (93% of its time in a profile: row scan 56%, FoldLog 32%, decode 11%), taking about 0.45 s on this repo against 0.02 s for an as-is read with a warm index, and 1.6 s on the 50k-node benchmark. Folding only the in-flight change sets a what-if needs over the stored as-is rows removes the scan; meeting the 50 ms budget for what-if reads would further need a SQL overlay reader and stays in the older proposed change set 01M5SWTFB3.

## title

Milestone 15: the feed shows the work, what-if reads stop rereading the log

## Links

### In

- derived-from [m15-picked-by-orchestrator](../decision/m15-picked-by-orchestrator.md)
- derived-from [feed-hides-actor-only](../decision/feed-hides-actor-only.md)
- derived-from [whatif-folds-inflight-over-asis](../decision/whatif-folds-inflight-over-asis.md)

## History

- `01M6BH7V2AFN6B2PYHQ4NFY3AG` Plan milestone 15: the feed shows the work, what-if reads stop rereading the log, picked by the orchestrator under the owner's standing instruction
