# m12-orchestrator-picks

source · `01M6ATAWYP6483ZDF7EDWKZKZP`

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

8cf978e2c4f4ae0d5d940d96bb2ad6991c6469421028e99540b769db120409dd

## summary

With milestones 10 and 11 queued, the orchestrator planned milestone 12 from what it still watched for by hand overnight. The disk filled (ENOSPC) and every agent session stalled for hours, and nothing in rhiza reads free space. run status prints a review check as 'failed after 2 cycles' while cycle 2 is open and its fixer working, because the label uses the last recorded outcome and counts the open cycle; the live app and run summary repeat the stale 'failed'. Failed and abandoned runs leave their worktree and cs/<slug> branch behind (16 cs branches here), since only close-out after a merge tears down. And the daemon picks milestone-10 and milestone-11 work before milestone-4 to 9, because sibling milestones sort by key as strings, though each carries a numeric milestone field. Research traced each: daemonPoll1 and the skip-reason table, runCheckInfo and Latest(), runRunAbandon and ship close-out, planCompareIntents.

## title

Milestone 12: the night runs clean

## Links

### In

- derived-from [m12-picked-by-orchestrator](../decision/m12-picked-by-orchestrator.md)
- derived-from [disk-guard-10-gib](../decision/disk-guard-10-gib.md)
- derived-from [open-cycle-reads-open](../decision/open-cycle-reads-open.md)
- derived-from [abandon-removes-clean-worktree](../decision/abandon-removes-clean-worktree.md)
- derived-from [milestones-by-number](../decision/milestones-by-number.md)

## History

- `01M6ATAR49P8K7TM7RCW7D20KF` Plan milestone 12: the night runs clean, picked by the orchestrator under the owner's standing instruction
