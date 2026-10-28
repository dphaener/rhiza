# ship-change-file-verified

evidence · `01M60APMY1KVHGHY90AM08AFXM` · actor [claude-code-4c6e2a4b](../actor/claude-code-4c6e2a4b.md)

## result

PASS

## statement

At 7c7f59d (diff 228dd44..HEAD: ship-change.toml, workflow_test.go, run-ship-change*.txtar) I ran: go test ./cmd/rhiza -count=1 -run 'TestScripts/(run-ship-change|describe-golden)' -v, all PASS. run-ship-change starts the built-in ship-change (no --workflow) with fake rex/claude/gh on an accepted change set of two leaves (second depends-on first); it checks builder then verifier for first, then second, two review check sessions, review/review pass then open-pr, one gh pr create whose body file is the review summary (# Review of run, ## Checks of step review/review), stopped with signoff_waiting and no pr merge and the change still accepted; after rhiza run approve it succeeds through rebase, land, merge (squash pinned to the landed head), main check, close-out (worktree and branches gone) and ship, the change set landed with PR Evidence. ship-change.toml has build and review workflow steps (review skip=[signoff]), open-pr, wait-pr, pr-green, signoff, then the old steps; grep finds no shipit and no --pr-open. go test ./internal/workflow passes, with TestShipChangeBuiltin failing on any shipit or --pr-open and checking the step order. go test ./cmd/rhiza -run TestScripts/describe-golden -update left git status clean, so the golden is current. make fmt-check and go test ./... -count=1 all pass. make self-check failed only on a stale render of this verifier session's own actor (claude-code-4c6e2a4b), which is not in the commit.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
