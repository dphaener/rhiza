# conductor-after-m35-inflight

decision · `01M79KKVSMDP3B655ETK8YQ7JC`

## question

When should the new milestone's change sets start building?

## statement

Milestone-36 builds after milestone-35's three in-flight builds land (review role prompts, plan role prompts, tiered verifier), whose prompts become the spawned sessions' briefs; milestone-35's other change sets run alongside.

## wording

After m35's in-flight builds (Recommended), picked by the owner

## Links

### Out

- derived-from [conductor-workflow-interview](../source/conductor-workflow-interview.md)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
