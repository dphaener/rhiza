# actor-stamp-reuses-session-verified

evidence · `01M6M7YVVYCG0G6S101MEP5V2H` · actor [claude-code-edb0a977](../actor/claude-code-edb0a977.md)

## result

PASS

## statement

At 02c4a07f, read the diff to recordAgentActor in internal/api/session.go: before writing, it scans the model's live actors (model.Children skips retired nodes) for one whose session field equals the session ID and returns that actor's handle and recording change set without writing anything. Ran go test -count=1 -v -run 'TestRecordAgentActor|TestSession' ./internal/api: TestRecordAgentActorReusesSession PASS. That test calls recordAgentActor twice for one session ID on one workspace, checks the same handle and change set come back with the change set count unchanged, then checks a different session ID gets a new handle and exactly one new change set. go vet ./internal/api clean, gofmt clean, full go test -count=1 ./internal/api ok (451s).

## Links

### Out

- verifies [actor-stamp-reuses-session](../capability/rhiza/runtime/actor-stamping/actor-stamp-reuses-session.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
