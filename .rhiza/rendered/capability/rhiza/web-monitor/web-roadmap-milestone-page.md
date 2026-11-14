# web-roadmap-milestone-page

capability · `01M7B38BDN8999D5B70GG5HASC` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run TestRoadmapMilestones shows, for a two-project root, no row from the other project, no finished milestone among the open ones, the groups in order with a change set whose run is going under building now and a proposed change set under proposed, titles from the change set's title field with the intent as the clamped sub-line, the state words, a progress bar whose segments and counts match the milestone's landed, accepted, proposed and withdrawn change sets, the done points inside a closed fold, the backlog and an umbrella group after the milestones, and the finished milestones in one in-page fold with no link to another page.

## statement

The roadmap page shows its project's open (not finished) milestones from the milestones operation, as in pages/roadmap.html: each with its state in words, a segmented progress bar and counts of landed, accepted, proposed and withdrawn change sets, its done points folded, and its change sets grouped (building now, from the runs; accepted; proposed; landed; drafts; withdrawn ones counted in the bar and the landed group's fold) with the short title, the intent clamped to two lines and a state pill; after the milestones, the backlog and each umbrella, from the roadmap operation as today; the finished milestones, and a group's change sets past those shown, fold into disclosures that expand in place on the roadmap, opening no other page.

## Links

### Out

- depends-on [web-roadmap-milestone-data](web-roadmap-milestone-data.md)
- decided-by [web-redesign-direction](../../../decision/web-redesign-direction.md)

### In

- verifies [web-roadmap-milestone-page-verified](../../../evidence/web-roadmap-milestone-page-verified.md)
- verifies [web-roadmap-milestone-page-qa-1](../../../evidence/web-roadmap-milestone-page-qa-1.md)
- verifies [web-roadmap-milestone-page-verified-2](../../../evidence/web-roadmap-milestone-page-verified-2.md)
- verifies [web-roadmap-milestone-page-qa-2](../../../evidence/web-roadmap-milestone-page-qa-2.md)
- verifies [web-roadmap-milestone-page-verified-3](../../../evidence/web-roadmap-milestone-page-verified-3.md)
- verifies [web-roadmap-milestone-page-qa-3](../../../evidence/web-roadmap-milestone-page-qa-3.md)
- verifies [pr-576-web-roadmap-and-change-set](../../../evidence/pr-576-web-roadmap-and-change-set.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
