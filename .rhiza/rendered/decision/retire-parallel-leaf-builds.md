# retire-parallel-leaf-builds

decision · `01M79KKVRZ6CQQWEV3YXNE4DYX`

## question

What should happen to milestone-35's 'parallel leaf builds with file locks' change set (01M76XTN), now that the conductor would decide parallelism?

## statement

Milestone-35's change set for daemon-driven parallel leaf builds with file locks (01M76XTN06GSM9) is withdrawn: the conductor decides which builders to spawn in parallel, and the spawn primitive enforces slots and load.

## wording

Retire it (Recommended), picked by the owner

## Links

### Out

- derived-from [conductor-workflow-interview](../source/conductor-workflow-interview.md)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
