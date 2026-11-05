# idle-run-session-stalls

decision · `01M6N9FFNSVCXW4GGGDGXERBZC`

## reasoning

internal/session/status.go Stalled returns false for every hook state but working, so a session that hit Stop (idle) while its background shell ran is invisible to runStalledAll in internal/api/run.go and to the inbox; that is the bug the change set's intent names. A run step waits only for `rhiza session done`, so idle without done means the run waits; the shell is not observable from hook status, so the rule keys on idle-without-done rather than on detecting a shell. Exclusions mirror the existing ones (usage-limit pause) plus states that already have an owner or arbiter path (Asked, blocked). Today nudging exists only as the inbox's manual n (internal/tui/inboxapp nudge sending carryOn via session.message), so keeping that path covers attended runs; the intent says the session should be nudged, and the owner's stated goal is that the daemon does all orchestration in autonomous mode, so the daemon nudges there. One nudge per stall avoids a loop of messages to a session that is genuinely stuck, leaving that case to the owner.

## statement

A live run-step session that has not said `rhiza session done` counts as stalled once its hook status has gone unchanged for more than StallAfter (20 min): working in any run, or idle in a session an autonomous run launched (session.Record.Autonomous), whether or not it holds a background shell. Never stalled: an idle session of an owner-attended run (it may wait on the owner, per idle-stall-autonomous-only), a session waiting on the arbiter (Status.Asked set), one paused at its usage limit, one blocked on a permission prompt (already its own inbox item), and an exited one. A stalled session shows in the inbox as today, where the owner nudges it with n. When the run itself is autonomous (the run's own Autonomous, decided per run by the autonomy config, not the daemon's --autonomous flag, which every shipped daemon carries), the daemon nudges it itself, once per stall (keyed on the report time it stalled at), with a message that tells it to check any background shell's result and carry on, or run `rhiza session done` if its work is finished; a session still stalled after that stays in the inbox for the owner and gets no repeat nudges. A session of an owner-attended run is never nudged by the daemon.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-f0ae9dff-1](../question/question-f0ae9dff-1.md)
- supersedes [stalled-after-20-minutes](stalled-after-20-minutes.md)
- derived-from [idle-stall-plan](../source/idle-stall-plan.md)

### In

- decided-by [run-session-idle-stall](../capability/rhiza/workflows/workflow-runs/run-status/run-session-idle-stall.md)
- decided-by [daemon-nudges-stalled](../capability/rhiza/orchestration/daemon/daemon-nudges-stalled.md)
- decided-by [run-status-stalled-session](../capability/rhiza/workflows/workflow-runs/run-status/run-status-stalled-session.md)
- decided-by [inbox-stalled-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
