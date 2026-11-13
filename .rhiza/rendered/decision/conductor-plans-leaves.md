# conductor-plans-leaves

decision · `01M79KKVQ03YYFAZDYWBBZW5H8`

## options

```json
[
  "Inside the conductor",
  "Separate plan run"
]
```

## question

Where should planning a change set's leaves happen, when an accepted change set has none?

## statement

A change set with no leaves gets them from its own conductor's first phase, not from a separate plan run.

## wording

Inside the conductor (Recommended), picked by the owner

## Links

### Out

- derived-from [conductor-workflow-interview](../source/conductor-workflow-interview.md)

### In

- decided-by [m36-conductor](../intent/milestone-36/m36-conductor.md)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
