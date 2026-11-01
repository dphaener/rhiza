# m19-orchestrator-picks

source · `01M6C3PXC1TTWTXQNSGCMX0MDZ`

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

6db92aaa4eced9c4c582bb284616cac5fde5a226f390c9796b11887e322161c5

## summary

While the queue ran one change set at a time, the orchestrator found that a plain rhiza ls in the main checkout took 0.28 s against 0.03 s for an explicit view. The checkout's current change set was a long-landed one, and openReader, when no view is given, looks up the effective change set's state with w.changeSet, which folds every change set in the log; a landed one then gives the as-is view anyway. The index already stores each change set's state (the change_sets table, read by xi.ChangeSets for view refs). Clearing the current change set by hand brought the read to 0.03 s. Every agent and app read without --view in a worktree with a current change set pays this fold.

## title

Milestone 19: default reads stay fast

## Links

### In

- derived-from [m19-picked-by-orchestrator](../decision/m19-picked-by-orchestrator.md)
- derived-from [default-view-state-from-index](../decision/default-view-state-from-index.md)

## History

- `01M6C3PTXZA57M17Y2Y9DX6NW9` Plan milestone 19: default reads stay fast, picked by the orchestrator under the owner's standing instruction
