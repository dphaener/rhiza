# spawn-done-plan

source · `01M7M2D9YCHMXKE0SJ0G18V686`

## quotes

```json
[]
```

## raw_sha256

2d47281ed31730b96cdcdcb3cfb96e48604ab2c11f80558526e982863c46ee8e

## session

66a967e2-de01-4388-a6b3-01a1e8241005

## summary

Autonomous plan run (no owner). Research: run.spawn (internal/api/run_spawn.go) waits in a loop on childOutcome, which ends only once the child's session record has an exit. A child that ran rhiza session done keeps its agent at the prompt, so with no daemon nothing ends it and spawn, and the conductor step waiting on it, hang until someone closes the child's pane. The engine already treats a done mark as finished for steps with until_when = done (readyForUntil, then closeSessionPane through stepRunner.ClosePane). The plan mirrors that in spawn: a done mark ends the wait with pass, and spawn closes the child's pane itself. childOutcome is shared with stepRunner.Child, which the engine's reconcileSpawns reads, so a restarted driver agrees with the spawn process. The change set serves m4-runs-hold-up, as the recent bug fixes on runs do. No questions were raised; no owner words exist.

## title

Plan of the spawn-waits-for-exit bug, autonomous

## Links

### In

- derived-from [spawn-ends-on-done-mark](../decision/spawn-ends-on-done-mark.md)
- derived-from [spawn-fix-stays-blocking](../decision/spawn-fix-stays-blocking.md)

## History

- `01M7M27J4MMS8X5PT4W3EXKDTM` Fix rhiza run spawn waits for the child agent's process to exit, and with no daemon nothing closes a child that already said rhiza session done, so spawn and its run hang until someone closes the child's pane
