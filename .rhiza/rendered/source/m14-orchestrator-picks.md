# m14-orchestrator-picks

source · `01M6BD9RFTEFB4RG2QFGXG6B22`

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

cf569a1111091d3ad167fe9ab2380de92a6acf7552cf4ccff205c1dc5d2ef75a

## summary

With milestone 13 the last one queued, the orchestrator planned milestone 14 under the owner's standing instruction. There is no search today: ls filters by kind only and no app filters, while the Area views umbrella still holds text-search ('genesis find searches the model's text'). The index uses modernc.org/sqlite with FTS5 compiled in, as-is and to-be views only; what-if views are built in memory. And the daemon's notifications reach only Rex (an OSC status report per new inbox item, recorded as notified only when Rex is attached), so an owner away from the screen hears nothing. Rhiza has no settings file; the daemon is configured by flags and the context already reads git config. The orchestrator carves text-search out of Area views (as the change feed was) and adds a notify hook under the landed notifications capability rather than carving Area hooks' broader model-event-hooks.

## title

Milestone 14: find it, hear about it

## Links

### In

- derived-from [m14-picked-by-orchestrator](../decision/m14-picked-by-orchestrator.md)
- derived-from [search-fts-asis-scan-whatif](../decision/search-fts-asis-scan-whatif.md)
- derived-from [notify-command-flag-or-git-config](../decision/notify-command-flag-or-git-config.md)

## History

- `01M6BD9FT5SGV1CMZXV25XZCMR` Plan milestone 14: find it, hear about it, picked by the orchestrator under the owner's standing instruction
