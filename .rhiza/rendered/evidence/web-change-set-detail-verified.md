# web-change-set-detail-verified

evidence · `01M7BKHJPJSV8MYXNA0CYNS97N` · actor [claude-code-b1e41eeb](../actor/claude-code-b1e41eeb.md)

## result

PASS

## statement

At a506d6b48 (attempt 2) re-ran the Checks: go test -trimpath -count=1 ./internal/web -run TestChangeSetPage (PASS, exit 0), rhiza ship generate (exit 0), rhiza render (exit 0); go test ./internal/web passes. Read the leaf's diffs (13ff93847, a506d6b48): only internal/web/change.go, change_test.go, templates/change.html and run records changed, no TODOs or stubs. The test now covers the acceptance on a landed change set: Shipped time, PR #498 link, run 'ship, 4m 30s', Serves done point, no change set ID or handles, three leaves with statement, acceptance and 'Verified, has evidence', decision/question/source linking to /p/proj/nodes/<id>, Proposed/Accepted/Landed history; a leafless draft says it has none yet. All three attempt-1 findings resolved.

## Links

### Out

- verifies [web-change-set-detail](../capability/rhiza/web-monitor/web-change-set-detail.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
