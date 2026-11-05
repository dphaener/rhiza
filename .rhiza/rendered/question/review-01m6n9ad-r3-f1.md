# review-01m6n9ad-r3-f1

question · `01M6NAB0JB7YTF7RJ68MP9DF2H`

## statement

Plan review round 3 left finding f1 (contradiction on run-session-idle-stall): The widened rule stalls any live run-step session idle >20 min without done, excluding only Status.Asked, usage-paused, blocked and exited. But in an owner-attended run an interactive step (e.g. a plan run's planner interviewing the owner, see rhiza-plan and the plan workflow) sits idle at its prompt waiting on the owner's reply, not on itself; Status.Asked is only set for autonomous sessions (internal/session/status.go), so such a session would be listed as 'stalled' in the inbox and run status, and the inbox's n would send it carryOn ('You have been silent...') when it is the owner who owes the next word. The Decision idle-run-session-stalls and this leaf's statement/acceptance should say how an attended session awaiting the owner is told apart (e.g. stall idle only in autonomous runs, or only for non-interactive steps), or state explicitly that it is accepted noise. The critic confirmed it: Verified: Status.Asked is set only when rec.Autonomous (internal/api/session_run.go:794-809), so an attended session's question to the owner never sets it. The plan workflow's interview step (internal/workflow/builtin/plan.toml) is an agent step that ends only on `rhiza session done` and spends most of its life idle at its prompt waiting for the owner's reply. Under run-session-idle-stall and idle-run-session-stalls, which exclude only Asked, usage-paused, blocked and exited, that session counts as stalled after 20 idle minutes. It would show in run status and the inbox, and the owner's n would send it carryOn, when it is the owner who owes the next word. The leaf and Decision must say how a session waiting on the owner is told apart, or accept that noise in so many words. How should it be settled?

## Links

### Out

- blocks [run-session-idle-stall](../capability/rhiza/workflows/workflow-runs/run-status/run-session-idle-stall.md)

### In

- answers [idle-stall-autonomous-only](../decision/idle-stall-autonomous-only.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
