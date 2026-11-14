# web-roadmap-milestone-page-qa-3

evidence · `01M7BY49HJFHSRES7CB9KZRRW9` · actor [claude-code-bc270773](../actor/claude-code-bc270773.md)

## result

PASS

## statement

QA at 08472f9c4, scratch-build daemon (paused, max-runs 0) on a two-project root: /roadmap redirected 302 to /p/rhiza/roadmap, which rendered in 1.5-5s (it was 50-91s at cycle 2). A script checked all 30 open milestones against /api/milestones: the kicker state, the bar aria-label counts of landed, accepted, proposed and withdrawn, the set of rows, each row's title and its intent sub-line, and that the only withdrawn row sits in the Landed fold ('8 more landed, 1 withdrawn'). Every check matched. The backlog ('Not in a milestone') and the umbrella came after the milestones, '6 finished milestones' was one in-page fold with no links, every link stayed under /p/rhiza, and every ID on the page came from the model's own text. /p/other/roadmap showed only its own milestones, with proposed (PR E), landed (LD B), drafts (DR D) and withdrawn rows in their groups. Three parallel loads gave byte-identical pages. Building now was not exercised (no runs in the clone). Rough edges: a milestone with only withdrawn change sets renders an empty Landed heading, and withdrawn rows show their intent as the title.

## Links

### Out

- verifies [web-roadmap-milestone-page](../capability/rhiza/web-monitor/web-roadmap-milestone-page.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
