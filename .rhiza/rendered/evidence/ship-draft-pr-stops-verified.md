# ship-draft-pr-stops-verified

evidence · `01M73J37SGGC19TNX946Y1849R` · actor [claude-code-87997a4c](../actor/claude-code-87997a4c.md)

## result

PASS

## statement

Verifier at 49db34e1: go test -trimpath ./internal/api -run 'TestShipCheck|TestShipMerge' -count=1 passes, including new TestShipCheckDraftStops (fake gh isDraft true on a green PR, plain and --settled: exit 75, stop: line naming #7 as a draft, Passed false, no checked head recorded) and TestShipMergeDraftStops (recorded head, pr view 7 --json isDraft,state reports draft: exit 75, stop line, no gh pr merge or pr ready call, state head kept, no merge commit); existing ready-PR green and merge tests still pass with the isDraft field added. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-merge-draft|describe-golden|run-ship-change)$' passes: run-ship-merge-draft turns the fake PR draft after rebased-merge-green, the run stops at merge with 'stop: pull request #7 is a draft', status shows current step merge, no merge visit 2 and no gate-rebase-merge, zero pr merge and no pr ready in gh.log; after removing the draft flag rhiza run resume succeeds, merges once and the change set lands. describe-golden was regenerated for the new ShipMergeResponse.stop field. Diff read: draft checked in ship check after the merged short-circuit and in ship merge before fetch/merge; Rhiza never calls gh pr ready.

## Links

### Out

- verifies [ship-draft-pr-stops](../capability/ship-through-rhiza/ship-draft-pr-stops.md)

## History

- `01M73G0WSCCBG7H96NDY1H4KQ6` Fix watchdog unmerged PR at step gate-rebase-merge: run 01M73CQ4D71CW038ZBNSQB1SJF; log: PR green at 0b28569 and unmerged after step gate-rebase-merge failed: step "gate-rebase-merge" went over max_visits 6
