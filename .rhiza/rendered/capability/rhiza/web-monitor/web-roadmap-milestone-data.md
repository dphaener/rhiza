# web-roadmap-milestone-data

capability · `01M7B38B8RHGSEEJ3KW7TGGNSJ` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/api -run TestMilestonesOperation returns, for a fixture with milestone A (a draft, a proposed, an accepted, a landed and a withdrawn change set; done point A1 served only by the landed one, done point A2 by the accepted one), A1 met and A2 not, every change set with its state, and A building; milestone B with only a proposed change set planned, milestone C with only drafts designing, milestone D whose one done point's only change set is landed finished, and a finished count of 1; the roadmap operation's response for the same fixture is unchanged (proposed and accepted only); go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes with the golden regenerated.

## statement

A new milestones operation (rhiza milestones --json), separate from the roadmap operation, which stays as it is, gives each of a project's milestones its done points (its child intents), each with whether it is met, its change sets in every state (draft, proposed, accepted, landed, withdrawn) with each one's state, and its state word; and counts how many of the project's milestones are finished. A done point is met when at least one change set serves it (directly or through a child intent) and every one of those that is not dropped (withdrawn, rejected or superseded) is landed. A milestone is finished when it has done points and all are met; else building when it has an accepted or a landed change set; else planned when it has a proposed one; else designing.

## Links

### Out

- decided-by [web-redesign-direction](../../../decision/web-redesign-direction.md)

### In

- depends-on [web-roadmap-milestone-page](web-roadmap-milestone-page.md)
- verifies [web-roadmap-milestone-data-verified](../../../evidence/web-roadmap-milestone-data-verified.md)
- verifies [web-roadmap-milestone-data-qa-1](../../../evidence/web-roadmap-milestone-data-qa-1.md)
- verifies [web-roadmap-milestone-data-verified-2](../../../evidence/web-roadmap-milestone-data-verified-2.md)
- verifies [web-roadmap-milestone-data-qa-2](../../../evidence/web-roadmap-milestone-data-qa-2.md)
- verifies [web-roadmap-milestone-data-qa-3](../../../evidence/web-roadmap-milestone-data-qa-3.md)
- verifies [pr-576-web-roadmap-and-change-set](../../../evidence/pr-576-web-roadmap-and-change-set.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
