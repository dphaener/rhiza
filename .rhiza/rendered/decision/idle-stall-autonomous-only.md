# idle-stall-autonomous-only

decision · `01M6NAFDX55HKYG63E4QPDR308`

## reasoning

The critic verified that Status.Asked is set only when rec.Autonomous (internal/api/session_run.go around line 795), so in attended runs no status marks a session waiting on the owner, and the plan workflow's interview step (internal/workflow/builtin/plan.toml) sits idle at its prompt for most of its life until `rhiza session done`. Telling interactive steps apart per step would need a new workflow field and still misses attended build steps that ask the owner something; accepting the noise would make the inbox cry stalled on every interview and invite a carryOn nudge at the owner's own turn. The session record already carries Autonomous (internal/session/session.go:37), so keying the idle rule on it is cheap, needs no new state, and matches where the change set's intent bites: the owner's goal is that autonomous runs move without them, and idle-run-session-stalls already confines the daemon's nudge to autonomous runs. In an autonomous run every question goes to the arbiter (Asked set, excluded), so an idle session without done there really is the run waiting on itself. The working-unchanged rule stays for all runs, so attended runs lose nothing they had.

## statement

The idle-without-done stall applies only to sessions an autonomous run launched (session.Record.Autonomous true). A session of an owner-attended run stalls only as before: working, without done, with its hook status unchanged for more than StallAfter (20 min); an idle attended session is never stalled, because it may be waiting on the owner's reply (a plan run's interview, or a question put to the owner) and nothing in its status tells that apart from a session stuck after its background shell. So run status and the inbox never list an attended session waiting at its prompt, and the owner's n never sends it carryOn for silence the owner owns. The idle case is a known gap in attended runs, accepted: the owner is present and sees the pane. Run-session-idle-stall's statement and acceptance and idle-run-session-stalls take this exclusion, and the acceptance adds an idle record without done and with Autonomous false that is not stalled at T+21m.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6n9ad-r3-f1](../question/review-01m6n9ad-r3-f1.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
