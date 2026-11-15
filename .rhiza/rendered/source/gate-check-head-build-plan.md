# gate-check-head-build-plan

source · `01M7EHP05A88WBBMGB8MVC5WJX`

## quotes

```json
[]
```

## raw_sha256

0edb4c3f8b62cc2c8a4a9ec84600220e5520e4fad872a22c169fd4804993cd0b

## session

ac842bb4-6ecb-48ab-b534-01a1dd191802

## summary

Autonomous plan run with no owner. The orchestrator's message for this run diagnosed the loop: the conductor default's gate step runs `rhiza check` through sh in the run worktree, which resolves to the installed rhiza, so a branch that adds a CLI command, flag or setting fails the generated-reference comparison against docs/reference that the branch's own build passes, and the gate hands back to conduct again and again. Its hint: run rhiza check (with the docs/reference checks) with the PR head's build ($SHIP_HEAD_RHIZA), as ship open-pr's pre-push check already does, falling back to the installed rhiza only when no head build exists, with a testscript where a branch adds a documented CLI flag and the gate passes, all in one small change set. The planner read internal/run/engine.go (gate), internal/api/run.go (stepRunner), internal/api/check.go, internal/api/ship_build.go (runViaHead, used by land and render) and ship_prepush.go, and settled the design itself; each choice is recorded as an assumed Decision. Nothing is left open.

## title

Fix the gate step: autonomous plan run

## Links

### In

- derived-from [gate-check-via-head](../decision/gate-check-via-head.md)
- derived-from [gate-check-head-fallback](../decision/gate-check-head-fallback.md)

## History

- `01M7EHHAG3CTRQCJCWSG2W7BTZ` Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN)
