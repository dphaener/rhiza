# web-roadmap-milestone-page-verified-2

evidence · `01M7BTEKTHW1FYXQY89A8N303M` · actor [claude-code-fa2e3d38](../actor/claude-code-fa2e3d38.md)

## result

PASS

## statement

Verified attempt 2 at 35b8358ce. Re-ran: go test -trimpath ./internal/web/ (exit 0); go test -trimpath ./internal/web -run TestRoadmapMilestones -v (acceptance, PASS, exit 0); go test -trimpath ./internal/api/ -run 'Milestone|ReadOnly|WebMilestones' (exit 0); make self-check in a clean detached worktree of HEAD (ok; in the run checkout it fails only on this verifier's own unrendered actor record). Did not re-run the 25m full ./internal/api run, which the commit records as a timeout under host load with no failing test; it is over the command budget and CI runs it. Diff read: milestones.go now folds served change sets not already folded (landed ones) so their title field is shown; roadmap_test.go adds a landed change set whose title differs from its intent and asserts title then intent sub-line. TestRoadmapMilestones covers every acceptance point: two-project isolation, no finished milestone among open ones, group order with building now and proposed, title/intent, state words, bar segments and counts (7/2/1/1), closed done-points fold, backlog and umbrella after milestones, finished milestones in an in-page fold with no link. Nothing outside the leaf changed besides rendered actor files; no stubs or TODOs.

## Links

### Out

- verifies [web-roadmap-milestone-page](../capability/rhiza/web-monitor/web-roadmap-milestone-page.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
