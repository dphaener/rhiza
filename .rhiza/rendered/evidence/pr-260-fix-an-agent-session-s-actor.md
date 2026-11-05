# pr-260-fix-an-agent-session-s-actor

evidence · `01M6METKM1KH4MRQFXBV2AVB8Y` · actor [dphaener](../actor/dphaener.md)

## head

acc61ca6071591d79e07fbd0957351ea9532000c

## pr

https://github.com/dphaener/rhiza/pull/260

## result

1 check passed: check SUCCESS

## statement

Pull request #260 (https://github.com/dphaener/rhiza/pull/260) passed CI at acc61ca

## status

verified

## Links

### Out

- verifies [actor-record-in-checkout](../capability/rhiza/runtime/actor-stamping/actor-record-in-checkout.md)
- verifies [actor-stamp-reuses-session](../capability/rhiza/runtime/actor-stamping/actor-stamp-reuses-session.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
