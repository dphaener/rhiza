# conductor-spawned-sessions

decision · `01M79KKVQKMABKSY7Q0ZZC4JZF`

## options

```json
[
  "All spawned sessions",
  "Spawned, cheap roles in-session"
]
```

## question

How should the conductor run its builders, verifiers and reviewers?

## quote_where

2026-11-13, on the proposal

## statement

Every builder, verifier, reviewer, critic and fixer is a separate session the daemon spawns on the conductor's request, journaled and visible, never an in-session sub-agent.

## wording

All spawned sessions (Recommended), picked by the owner

## Links

### Out

- derived-from [conductor-workflow-interview](../source/conductor-workflow-interview.md)

### In

- decided-by [m36-visible](../intent/milestone-36/m36-visible.md)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
