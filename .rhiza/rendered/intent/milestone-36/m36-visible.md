# m36-visible

intent · `01M79KKVM10T0D703RCS27QB7M` · parent [milestone-36](../milestone-36.md)

## statement

Every session a run spawns is visible as a child of the run in rhiza run status and the web run page, with its role, leaf, model and outcome, and can be opened or resumed like any other session.

## Links

### Out

- decided-by [conductor-spawned-sessions](../../decision/conductor-spawned-sessions.md)

### In

- serves change set `01M79KGRYE538GJC5NA2VNYG4K` (The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
