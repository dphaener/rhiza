# conductor-three-checks

decision · `01M79KKVR7T0Z38NYD6T1TTWGR`

## options

```json
[
  "Three: Plan, Code, QA",
  "Two: Plan and Code",
  "Keep all five"
]
```

## question

How many review checks should the default conductor run?

## statement

The default conductor runs three review checks: Plan (model and compliance), Code (quality) and QA (the built binary does what the leaf says); the scope check is dropped; a critic runs only on a check with a blocker or major finding.

## wording

Three: Plan, Code, QA (Recommended), picked by the owner

## Links

### Out

- derived-from [conductor-workflow-interview](../source/conductor-workflow-interview.md)

### In

- decided-by [m36-conductor](../intent/milestone-36/m36-conductor.md)
- decided-by [review-three-checks](../capability/rhiza/workflows/built-in-workflows/review-three-checks.md)
- decided-by [critic-on-major](../capability/rhiza/workflows/built-in-workflows/critic-on-major.md)
- decided-by [check-plan-code-topics](../capability/rhiza/guide/agent-guide/check-plan-code-topics.md)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
