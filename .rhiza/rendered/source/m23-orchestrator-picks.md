# m23-orchestrator-picks

source · `01M6D12Q1PQ43BC4YD0Z4SQMWT`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I can interact with rhiza via the mini TUI apps to examine the full health of the application, start new changes, report bugs, start new features, etc, without having this interactive session running.",
    "where": "orchestrator session, 2026-11-01"
  }
]
```

## raw_sha256

a0c7b23ae888009762d7101891a4b9df89e0bec3662227796b8f46a84d431ce6

## summary

At 2026-11-02 02:20Z the live daemon settings run stopped at wait-pr on a code conflict with main, as the new wait step does. Resolving it by hand, the orchestrator found that rhiza question answer refused the run's rebase-conflict question as 'not an open question blocking this change set': shipRaiseQuestion (internal/api/ship_rebase.go) adds the question with only a statement, no blocks link to the change set, so the blocker rule question.answer and the inbox use cannot see it. The orchestrator answered with a hand-written Decision instead. For the owner to answer a conflict from the inbox app, the question must block its change set.

## title

Milestone 23: answer a conflict from the inbox

## Links

### In

- derived-from [m23-picked-by-orchestrator](../decision/m23-picked-by-orchestrator.md)
- derived-from [conflict-question-blocks](../decision/conflict-question-blocks.md)

## History

- `01M6D12MQKVB82JNS8JR91K4FC` Plan milestone 23: answer a conflict from the inbox
