# actor-stamp-reuses-session-qa-1

evidence · `01M6M851NMHM3XRM1H3Z233WT6` · actor [claude-code-7e6aebbb](../actor/claude-code-7e6aebbb.md)

## result

PASS

## statement

QA at ca2af093 by claude-code-7e6aebbb: the CLI gives each session a fresh ID, so the same-session repeat is only reachable in-process; go test ./internal/api -run TestRecordAgentActorReusesSession -count=1 passed (second call returns the same handle and change set and writes nothing, a different session gets a new actor). In the hands-on build run with the built binary, the run worktree's log held exactly one 'Record agent actor' change set per session (builder and verifier), with no duplicates after the verifier's checkout copy merged.

## Links

### Out

- verifies [actor-stamp-reuses-session](../capability/rhiza/runtime/actor-stamping/actor-stamp-reuses-session.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
