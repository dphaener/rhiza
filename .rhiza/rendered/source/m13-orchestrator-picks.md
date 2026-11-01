# m13-orchestrator-picks

source · `01M6B3NWCTG92ZGHQE28XWAF1V`

## quotes

```json
[
  {
    "by": "owner",
    "text": "if you run out of work keep picking high value tasks. don't stop and wait for my input.",
    "where": "orchestrator session, 2026-10-31"
  },
  {
    "by": "owner",
    "text": "find the fastest possible route to this being usable, by me, from rex, using the TUI to launch things",
    "where": "orchestrator session, 2026-10-30"
  }
]
```

## raw_sha256

eb98c2d1dec5118299136d83e59246b8b384aef9d94dc1729724cdcbb7b7bba2

## summary

With milestone 12 the last one queued, the orchestrator planned milestone 13 under the owner's standing instruction, aimed at the owner coming back to a night of unattended work: nothing lists landed change sets by time (change ls sorts by ID; land times only show in change show's history), so catching up means reading git log and PRs. The founding conversation decided there is a change feed (events --since and --follow) with cursor-paged lists; it sits in the Area views umbrella change set. Research found land entries carry ULIDs and times, the index records each change set's terminal entry, and ship.land writes PR Evidence (pr field) that survives close-out, so a feed can show PR numbers. Because a change set lands on its branch before it merges, a land-entry cursor can skip work merged later, so the app remembers which change sets the viewer has seen rather than a cursor. Separately, TestReadBudget fails under load from outside the test run (other worktrees, agents), with no retry.

## title

Milestone 13: catch up at a glance

## Links

### In

- derived-from [m13-picked-by-orchestrator](../decision/m13-picked-by-orchestrator.md)
- derived-from [feed-lists-landed-by-land-entry](../decision/feed-lists-landed-by-land-entry.md)
- derived-from [changes-app-remembers-seen](../decision/changes-app-remembers-seen.md)
- derived-from [budget-retries-over-reads](../decision/budget-retries-over-reads.md)

## History

- `01M6B3NBFS2XB29R29457RSCGA` Plan milestone 13: catch up at a glance, picked by the orchestrator under the owner's standing instruction
