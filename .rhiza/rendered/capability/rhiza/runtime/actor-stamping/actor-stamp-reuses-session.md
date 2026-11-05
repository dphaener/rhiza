# actor-stamp-reuses-session

capability · `01M6M2ZS1VT4ZWAMBM148Z3RG6` · parent [actor-stamping](../actor-stamping.md)

## acceptance

A Go test calls recordAgentActor twice for one session ID on one workspace: the second call returns the first call's handle and writes no new change set; a different session ID still gets its own new actor

## statement

Recording an agent session's actor reuses an existing, unretired actor whose session field matches that session instead of adding a new one, so a repeated stamp writes nothing and returns the same handle

## Links

### Out

- decided-by [actor-record-once-per-run](../../../../decision/actor-record-once-per-run.md)

### In

- verifies [actor-stamp-reuses-session-verified](../../../../evidence/actor-stamp-reuses-session-verified.md)
- verifies [actor-stamp-reuses-session-qa-1](../../../../evidence/actor-stamp-reuses-session-qa-1.md)
- verifies [pr-260-fix-an-agent-session-s-actor](../../../../evidence/pr-260-fix-an-agent-session-s-actor.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
