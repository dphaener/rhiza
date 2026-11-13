# m36-conductor

intent · `01M79KKVJZBSHJ5RXCG4JKBSP4` · parent [milestone-36](../milestone-36.md)

## statement

One conductor session conducts a change set end to end: it plans the leaves when there are none, spawns a builder and then a verifier per leaf, spawns the Plan, Code and QA review checks with a critic only on a blocker or major finding, spawns a fixer for confirmed findings, and lands the change set on the branch.

## Links

### Out

- decided-by [conductor-plans-leaves](../../decision/conductor-plans-leaves.md)
- decided-by [conductor-three-checks](../../decision/conductor-three-checks.md)

### In

- serves change set `01M79KGZ51V2F5YYWMMQTGR6NW` (A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands)
- serves change set `01M79KH59D0HEME6B85PXXB056` (The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
