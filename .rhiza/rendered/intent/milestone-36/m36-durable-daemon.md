# m36-durable-daemon

intent · `01M79KKVKEPHJQ6MH0SZ58FG2Z` · parent [milestone-36](../milestone-36.md)

## statement

The daemon is dumb but durable: workflows are built from the agent, spawn, gate, push/pr, wait, merge, signoff and hook primitives; any red result hands back to a named agent step with the cause, at most three times before the run parks with an inbox item; every step and spawned session is journaled, and killing the daemon at any point loses nothing.

## Links

### In

- serves change set `01M79KGRYE538GJC5NA2VNYG4K` (The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
