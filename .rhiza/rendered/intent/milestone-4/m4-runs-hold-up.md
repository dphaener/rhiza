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
- serves change set `01M732D2EZSS20MWDHZGG66BJ5` (Fix A run whose change bumps the log entry format can never pass review/code. Run 01M72YSG33J9JB66SYG9AD0Z0F (Docs sign-off workflow, change set 01M71G46Q5378CYYVEHFF2EMP3) failed at review/code with three causes like: 'rhiza check: fail log .rhiza/log/01M71G4BKFHP1E172RR6H82WFT/01M731FHDC43E8XR3XAPNAFE1Y.json: entry: newer format: field "v": 5, this Rhiza reads 4'. The builder's own binary wrote v5 entries (pinning a change set's workflow), but run code-check runs rhiza check in-process with the installed daemon binary (internal/api/run_codecheck.go runCheck; entry.ErrNewerFormat at internal/entry/entry.go:43), which reads only v4. #417 (01M6ZT7P) skipped the reference rules for the same root cause but not this one. Expected: code-check checks a run's worktree with that worktree's own build (or otherwise accepts the format its change introduces), so format bumps can ship through the daemon. Unstick: once fixed and merged (daemon self-updates), rhiza run resume --retry 01M72YSG33J9JB66SYG9AD0Z0F.)
- serves change set `01M73FYE2PAG9E034MPBTTJ12D` (Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.)
- serves change set `01M73G0WSCCBG7H96NDY1H4KQ6` (Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6)

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
