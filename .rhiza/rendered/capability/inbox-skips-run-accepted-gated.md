# inbox-skips-run-accepted-gated

capability · `01M7ANT7A4WT45GF5PDBXQ0C57`

## acceptance

A new test in internal/api/inbox_run_accepted_test.go feeds inboxItems a plan-run journal with a gated entry for two change sets and an accepted entry for one of them, with changeState reporting both proposed: exactly one InboxAccept item is listed, for the change set the run did not accept; with no accepted entry both are listed. go test -trimpath ./internal/api/ -run 'TestInbox' and go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-owner-run-accept|stage-autonomy)' pass. Touches internal/api/inbox.go, internal/api/inbox_run_accepted_test.go

## statement

rhiza inbox raises no accept item for a change set a plan run gated when that same run journaled an accepted entry for it (st.Accepts, Verb accepted, Change the gated ID), even while main's model still says proposed; a gated change set the run did not accept (rhiza.autonomy.accept owner, or the review did not pass) keeps its 'waiting for acceptance' item with next rhiza accept --rule owner --change <id>

## Links

### Out

- decided-by [run-accepted-no-inbox-item](../decision/run-accepted-no-inbox-item.md)

### In

- verifies [inbox-skips-run-accepted-gated-verified](../evidence/inbox-skips-run-accepted-gated-verified.md)
- verifies [inbox-skips-run-accepted-gated-qa-1](../evidence/inbox-skips-run-accepted-gated-qa-1.md)
- verifies [pr-557-fix-the-inbox-asks-the-owner](../evidence/pr-557-fix-the-inbox-asks-the-owner.md)

## History

- `01M7ANMWSVKNBSNPYCZXPK41GC` Fix The inbox asks the owner to accept a change set that an autonomous plan run already accepted. On 2026-11-13, plan runs 01M7AH6Q9G (change set 01M79KGRYE) and 01M7AH6RHK (01M79KGHYF) journaled 'gated' then 'accepted' at 14:20Z (log entries on the plan branch, e.g. 01M7AMGNTAEAQSSDF5YFMHXPYD), opened plan PRs #550/#551 and moved on to wait-rebased; meanwhile rhiza inbox listed both as kind 'accept', 'waiting for acceptance', next 'rhiza accept --rule owner --change ...' (items e512e92, 400de15), because main's model still says proposed until the plan PR merges. The same happened for 01M79KGZ51 (run 01M7AH6R61) before its PR merged. Nothing needed doing: the items clear when the PR merges. Wanted: the inbox does not raise an accept item for a change set whose gated run journaled 'accepted' (or labels it 'accepted on branch, PR #N merging'), so the owner is never asked to do what the run already did.
