# plan-tiled-panes-autonomous

source · `01M651KJPSKDQJED0C7BS3ASMN`

## quotes

```json
[
  {
    "by": "owner",
    "text": "sounds like a rhiza problem. we should be co-locating like processes in the same window with tiled panes.",
    "where": "this conversation, after the plan runs failed on the window limit"
  }
]
```

## raw_sha256

40ab7acaeda288d2d1a903b69551f362525310eb7f20382d6a90d6c30857b4f0

## session

88b733b2-b1ce-4027-b01e-01a18a18b610

## summary

Autonomous plan run for change set 01M6516H (tiled run sessions). No owner was present, so no question was put to the owner. The plan rests on the owner's earlier decision that a run's processes share one window as tiled panes, recorded in window-limit-owner-words. Three leaves follow from it: a run's sessions open as panes of its own window; a session whose anchor pane is gone splits a live pane of that window instead of opening a window; review, critic and fixer sessions anchor to the run's pane. Open for the arbiter: the intent the change set serves, which link points at a node no longer live; the window a run's panes belong to; and whether an anchor that is gone should be an error rather than a split.

## title

Tiled run sessions: autonomous plan

## Links

### In

- derived-from [co-locate-run-panes](../decision/co-locate-run-panes.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
