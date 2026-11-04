# question-7346bd67-1

question · `01M6M2WCPY1DWW4GKJ8V8KAP0V`

## session

7346bd67-434f-4dba-9398-01a1a82c3dce

## statement

Which intent should change set 01M6M2R7 (duplicate agent actor records) serve? m4-runs-hold-up says runs finish 'with no swept-up actor records', and the earlier actor-record fix 01M602TH serves milestone-4. / How should a session's actor record stop being duplicated? Cause found in the code: session start (recordAgentActor, internal/api/session.go:668) writes the record only in the run worktree. A verifier works in its own checkout (dir), which lacks the record, so the agent hand-adds `add actor <handle> --record` there, and both get committed.

## Links

### In

- answers [actor-record-once-per-run](../decision/actor-record-once-per-run.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
