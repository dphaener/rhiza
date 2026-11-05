# m4-runs-hold-up

intent · `01M61R4RRY5ZEVDSZ0E2M11DRB` · parent [milestone-4](../milestone-4.md)

## statement

A build, review or ship-change run finishes without hand repair: no hand-made render commit, no second run on one change set, no swept-up actor records.

## Links

### Out

- decided-by [answer-8cd68bc7-6](../../decision/answer-8cd68bc7-6.md)

### In

- serves change set `01M61GW7K1PDG7HMFXMN3YNA7B` (Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.)
- serves change set `01M6J965BRJ2ZTDT7Y311PM96G` (Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless)
- serves change set `01M6M2R7RJCMNMMKH0KS2364BK` (Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records)
- serves change set `01M6M4YHKJAER13CM0YHZQFKCZ` (Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone)

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
