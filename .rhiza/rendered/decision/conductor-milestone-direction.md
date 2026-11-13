# conductor-milestone-direction

decision · `01M79KKVPGCBYHVX9H9WB8P5SM`

## quote

we also need to keep in mind that people should be able to bring their own workflows. the daemon should be dumb but durable. the original ask for orchestrating from the daemon was determinism but we're not getting that. the Damon should provide useful hooks for things that agents can't do, or shouldn't do. it should be as flexible as we can make it, but useful.

## statement

Rhiza's runs move from a daemon-driven workflow of about 56 steps to a conductor agent per change set over a dumb, durable daemon of a few primitives, with a gate-checked agent contract so teams bring their own skills, agents and workflows; determinism means durable decisions and deterministic gates, not scripting each agent judgment.

## Links

### Out

- derived-from [conductor-workflow-interview](../source/conductor-workflow-interview.md)

### In

- decided-by [milestone-36](../intent/milestone-36.md)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
