# window-limit-owner-words

source · `01M6516M9J4XRQADVHQ4DFK0MV`

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

5348e534f3c3579258e1c8669fc5e86695fa3189f51a5e1c35e7f935ad07f9a5

## summary

A plan run's review, critic and fixer sessions each open their own Rex window, so one session hit Rex's 64-window limit and failed four runs. The owner wants processes that belong to one run co-located in one window as tiled panes.

## title

Rhiza runs open too many Rex windows

## Links

### In

- derived-from [co-locate-run-panes](../decision/co-locate-run-panes.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
