# spawn-no-wait-plan

source · `01M7QE17XRJ2997XYN2BQQX1CQ`

## quotes

```json
[]
```

## raw_sha256

3144ab01f26be359b0f00b3ae7377e866bd0eea556cd157dec907f84fe778f6f

## session

dc492322-38b6-41d5-bece-01a1eedd8716

## summary

The conductor of run 01M7QDTMXTQE8QZQQ2PH2W22VB planned this change set with no owner present. Research: rhiza run spawn (internal/api/run_spawn.go) starts a child, journals a spawn entry and polls the child's record until it says done or exits, then journals spawn-ended and closes a said-done child's pane. The run's engine already polls a conductor step's children (reconcileSpawns records ends through ChildWatcher; askChildren has the arbiter answer their questions). Three design questions went to the arbiter, which answered with Decision spawn-no-wait-wake: --no-wait is opt-in, the engine sends the wake-up when reconcileSpawns journals a no-wait child's end, and a child's question wakes no one. Planned leaves: the --no-wait flag, the engine's wake-up (with the said-done child's pane closed, as the waiting spawn does), and the conduct guides telling a conductor to spawn with --no-wait and end its turn.

## title

Autonomous plan of change set 01M7M2FB (run spawn can return)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
