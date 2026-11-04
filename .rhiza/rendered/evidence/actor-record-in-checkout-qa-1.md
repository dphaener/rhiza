# actor-record-in-checkout-qa-1

evidence · `01M6M84SQSGJKGKBVWWR1B38RX` · actor [claude-code-7e6aebbb](../actor/claude-code-7e6aebbb.md)

## result

PASS

## statement

QA at ca2af093 by claude-code-7e6aebbb with the built binary (go build ./cmd/rhiza) on a scratch repo, fake rex/claude only standing in for the terminal and agent: rhiza run start --foreground --workflow build ran a builder and a verifier for leaf good (state succeeded). In the verifier's own checkout (.git/rhiza/checkouts/<run>-good-1) the log held its actor record change set 01M6M8339X..., diff -r against the run worktree's copy showed no difference, rhiza show claude-code-687aa7ac resolved, and rhiza add actor claude-code-687aa7ac --record was refused (exit 3, 'handle ... is taken by node'). I then committed the checkout's copy on a branch from the leaf commit and merged it into the run branch, which also committed the record: the merge was clean, rhiza render + rhiza check printed ok (exit 0), and the log held one record change set for that handle.

## Links

### Out

- verifies [actor-record-in-checkout](../capability/rhiza/runtime/actor-stamping/actor-record-in-checkout.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
