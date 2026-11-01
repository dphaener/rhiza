# m18-orchestrator-picks

source · `01M6BZTRS34CKEPT1JT06CRN8C`

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

f8e2b2d7590e4e09be3297f1fae03e0d5d54c278db584656c61038819ec5c930

## summary

After the owner asked what was left to build and said to keep going, the orchestrator planned milestone 18 from the small contained items left in the Area umbrellas. hook-commands (hook ls and hook run) finishes the hooks area begun by the notify hook and the lifecycle event hooks; it can be carved before the event hooks land because its change set depends on theirs, so the carve's what-if view includes its parent. Today hook ls could not see a running daemon's flags, since the daemon's state keeps only which items it handled, so the daemon records its commands and their sources. Notes on a change (decision notes-are-discussion) need storage: a new note node kind with a discusses link to anything is the smallest fit (adding the anchor kind touched four Go files), whereas reusing relates-to would blur its meaning and a new entry type would change folding; ops are not link endpoints, so a note on an op links to its change set and names the op. Actor stamping, also listed, is already landed (actor-stamping and launch-env-stamping), so it is dropped.

## title

Milestone 18: test your hooks, discuss a change

## Links

### In

- derived-from [m18-picked-by-orchestrator](../decision/m18-picked-by-orchestrator.md)
- derived-from [hook-ls-reads-daemon-state](../decision/hook-ls-reads-daemon-state.md)
- derived-from [notes-are-note-nodes](../decision/notes-are-note-nodes.md)

## History

- `01M6BVT9HNN801JPDD99CM049F` Plan milestone 18: test your hooks, discuss a change, picked by the orchestrator under the owner's standing instruction
