# m36-touches

intent · `01M79KKVNWQJPFDDBJVF43HG00` · parent [milestone-36](../milestone-36.md)

## statement

Over the same 30 runs, at most 3 need a person beyond the owner's own sign-offs and real questions: no hand-driving a ship, killing a process, nudging a session or retrying a step by hand.

## Links

### Out

- decided-by [conductor-targets](../../decision/conductor-targets.md)

### In

- serves change set `01M79KHRHFQYBAW375MA5D7V70` (Measure 30 runs on the conductor default against today's Rhiza and afk: median start-to-merge with planning, runs needing a person, tokens per change)
- serves change set `01M7MPAHJZY3RV8ANGGB660J52` (Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
