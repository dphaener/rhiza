# inbox-skips-run-accepted-gated-qa-1

evidence · `01M7AV25RPB7YDM74JMCY63SBN` · actor [claude-code-9422ca93](../actor/claude-code-9422ca93.md)

## result

PASS

## statement

QA cycle 1: built go build -trimpath ./cmd/rhiza to a scratch path (and the merge-base build for comparison). In a scratch Rhiza root/repo, a hold-workflow plan run with rhiza.autonomy.accept auto gated two change sets and journaled accepted for both on a plan branch; main got only their proposed log entries (change ls: proposed). rhiza inbox --json: base binary listed two accept items 'waiting for acceptance' (the bug); the built binary listed none. Dropping one accepted line from the run journal: the built binary listed exactly one accept item, for that change set, next rhiza accept --rule owner --change <id>, in --json and plain text. With rhiza.autonomy.accept owner (gated, no accepted entry): both accept items listed; after rhiza accept --rule owner on one, only the other remains. Scoped go test ./internal/api -run TestInbox and TestScripts/(inbox-owner-run-accept|stage-autonomy) both pass (ok).

## Links

### Out

- verifies [inbox-skips-run-accepted-gated](../capability/inbox-skips-run-accepted-gated.md)

## History

- `01M7ANMWSVKNBSNPYCZXPK41GC` Fix The inbox asks the owner to accept a change set that an autonomous plan run already accepted. On 2026-11-13, plan runs 01M7AH6Q9G (change set 01M79KGRYE) and 01M7AH6RHK (01M79KGHYF) journaled 'gated' then 'accepted' at 14:20Z (log entries on the plan branch, e.g. 01M7AMGNTAEAQSSDF5YFMHXPYD), opened plan PRs #550/#551 and moved on to wait-rebased; meanwhile rhiza inbox listed both as kind 'accept', 'waiting for acceptance', next 'rhiza accept --rule owner --change ...' (items e512e92, 400de15), because main's model still says proposed until the plan PR merges. The same happened for 01M79KGZ51 (run 01M7AH6R61) before its PR merged. Nothing needed doing: the items clear when the PR merges. Wanted: the inbox does not raise an accept item for a change set whose gated run journaled 'accepted' (or labels it 'accepted on branch, PR #N merging'), so the owner is never asked to do what the run already did.
