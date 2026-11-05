# actor-record-in-checkout

capability · `01M6M2ZMM64QSSKQTP1FKH76NJ` · parent [actor-stamping](../actor-stamping.md)

## acceptance

A testscript starts a build run's verifier session with its own checkout: the checkout's .rhiza/log holds the session's actor record change set with files byte-identical to the run worktree's, rhiza show <actor handle> resolves in the checkout, rhiza add actor <that handle> --record there is refused as a taken handle, and after both copies are committed and merged rhiza check passes with one actor node for the handle

## statement

When an agent session has its own checkout (a verifier's dir), session start copies the session's actor record change set, the same log entries with the same IDs, from the run worktree into that checkout, so the actor resolves there and committing it from both places yields one record

## Links

### Out

- decided-by [actor-record-once-per-run](../../../../decision/actor-record-once-per-run.md)

### In

- verifies [actor-record-in-checkout-verified](../../../../evidence/actor-record-in-checkout-verified.md)
- verifies [actor-record-in-checkout-qa-1](../../../../evidence/actor-record-in-checkout-qa-1.md)
- verifies [actor-record-in-checkout-verified-2](../../../../evidence/actor-record-in-checkout-verified-2.md)
- verifies [pr-260-fix-an-agent-session-s-actor](../../../../evidence/pr-260-fix-an-agent-session-s-actor.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
