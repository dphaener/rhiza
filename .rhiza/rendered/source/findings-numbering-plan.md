# findings-numbering-plan

source · `01M64JT7B39FATHKB4T1Z60KA0`

## quotes

```json
[]
```

## raw_sha256

8b2398a9d3bd72abef87d0d93214b1270bc927a7b4cc53e36c16eeacdd2fc0b6

## session

7cf182b9-2039-4d75-99d2-01a1892c6c98

## summary

Autonomous plan run, no owner present, so no owner words are quoted. The change set's intent: rhiza run findings numbers each finding f1, f2, ... from a count read in the run, so two concurrent calls in one session can journal the same ID. Planned outcome: findings are numbered under the run's lock, for the round, --check and leaf findings alike. Decisions below are the planner's assumptions, not owner decisions; the owner has not confirmed them.

## title

Findings numbering under the run's lock

## Links

### In

- derived-from [findings-under-run-lock](../decision/findings-under-run-lock.md)

## History

- `01M5ZNV40RB8VYBB0A3CPXCNXZ` rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).
