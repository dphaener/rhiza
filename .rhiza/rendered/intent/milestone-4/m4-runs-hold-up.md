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

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
