# web-change-set-detail

capability · `01M79Q8MWE8XTM6BXCY3EXST92` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run TestChangeSetPage shows a landed change set's facts and no change set ID, three leaves with statement, acceptance and a verified Evidence line, its decision, question and source linking to /p/<project>/nodes/<ref>, and its history; a draft with no leaves says it has none yet.

## statement

A change set's page shows, as in pages/change-set.html: its title, state and full intent; facts (when it shipped and its PR, its run and duration, the done point it serves), with no change set ID; its leaves with statement, acceptance and Evidence state; why it exists (decisions, questions, sources), each linking to its model page; and its history.

## Links

### Out

- decided-by [web-redesign-direction](../../../decision/web-redesign-direction.md)

### In

- verifies [web-change-set-detail-verified](../../../evidence/web-change-set-detail-verified.md)
- verifies [web-change-set-detail-qa-1](../../../evidence/web-change-set-detail-qa-1.md)
- verifies [web-change-set-detail-qa-2](../../../evidence/web-change-set-detail-qa-2.md)
- verifies [web-change-set-detail-qa-3](../../../evidence/web-change-set-detail-qa-3.md)
- verifies [pr-576-web-roadmap-and-change-set](../../../evidence/pr-576-web-roadmap-and-change-set.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
