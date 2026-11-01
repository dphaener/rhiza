# m21-orchestrator-picks

source · `01M6CA9SBSR9192WWQJY8MMCV2`

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

0d0932c6e97951ba4eecc5a4ba56e9632895118e8d286c44f38a9800b447ad1d

## summary

Two rough edges the orchestrator hit on 2026-11-01. During a daemon re-exec, rhiza daemon status printed only the process, with empty start and poll times, max runs 0 and no live runs, for several seconds: the new image takes the lock with a new start time, so status finds the state file names the old image and shows only the process, and the new image writes no state until its first poll ends (planning and notify take seconds), although the State comment says the file is written at start. And the main checkout kept a change set that landed on Oct 28 as its current one, because land clears the current change set only in the checkout that runs it, and work lands from run worktrees; writes and the model already treat a terminal current change set as none, but rhiza status and the session hook show it with no next step.

## title

Milestone 21: status tells the truth

## Links

### In

- derived-from [m21-picked-by-orchestrator](../decision/m21-picked-by-orchestrator.md)
- derived-from [daemon-writes-state-at-start](../decision/daemon-writes-state-at-start.md)
- derived-from [terminal-current-hint](../decision/terminal-current-hint.md)

## History

- `01M6CA9PFRY55ZTXNHWDM95Y9A` Plan milestone 21: status tells the truth, picked by the orchestrator under the owner's standing instruction
