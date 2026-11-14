# inbox-skips-run-accepted-gated-verified

evidence · `01M7ATJJ88AEGDJDNF3195GCZ5` · actor [claude-code-2e21957a](../actor/claude-code-2e21957a.md)

## result

PASS

## statement

At ce1602eb2 re-ran the Checks block: go test -trimpath ./internal/api/ -run TestInbox (exit 0, TestInboxSkipsRunAcceptedGated PASS), go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-owner-run-accept|stage-autonomy)' (exit 0, inbox-owner-run-accept and stage-autonomy PASS), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0, no tracked changes). Diff touches only internal/api/inbox.go (runAccepted skips a gated change set with an accepted Accepts entry before the proposed check), the new internal/api/inbox_run_accepted_test.go, and the builder's actor record. The test feeds inboxItems a plan-run journal gating C1 and C2 with changeState proposed: both listed with no accepted entry, only C2 (next ownerAcceptNext) with C1 accepted, both with C1 rejected. No stubs or TODOs; gofmt clean.

## Links

### Out

- verifies [inbox-skips-run-accepted-gated](../capability/inbox-skips-run-accepted-gated.md)

## History

- `01M7ANMWSVKNBSNPYCZXPK41GC` Fix The inbox asks the owner to accept a change set that an autonomous plan run already accepted. On 2026-11-13, plan runs 01M7AH6Q9G (change set 01M79KGRYE) and 01M7AH6RHK (01M79KGHYF) journaled 'gated' then 'accepted' at 14:20Z (log entries on the plan branch, e.g. 01M7AMGNTAEAQSSDF5YFMHXPYD), opened plan PRs #550/#551 and moved on to wait-rebased; meanwhile rhiza inbox listed both as kind 'accept', 'waiting for acceptance', next 'rhiza accept --rule owner --change ...' (items e512e92, 400de15), because main's model still says proposed until the plan PR merges. The same happened for 01M79KGZ51 (run 01M7AH6R61) before its PR merged. Nothing needed doing: the items clear when the PR merges. Wanted: the inbox does not raise an accept item for a change set whose gated run journaled 'accepted' (or labels it 'accepted on branch, PR #N merging'), so the owner is never asked to do what the run already did.
