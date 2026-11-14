# web-roadmap-milestone-page-qa-1

evidence · `01M7BN41B45DQEFX8ED73J40JN` · actor [claude-code-3866e93d](../actor/claude-code-3866e93d.md)

## result

FAIL

## statement

QA at 33cc5a571, a scratch-build daemon (paused) on a two-project root: /p/rhiza/roadmap showed 30 open milestones each with 'Milestone · <state>', a segmented bar and legend (e.g. 13 landed, 2 accepted, 1 withdrawn), done points in a closed fold, groups Accepted waiting, Proposed, Landed, Drafts, with folds such as '8 more landed, 1 withdrawn', then the backlog and the umbrella, and '6 finished milestones' in one in-page fold with no links; every link stayed under /p/rhiza and no ID was rendered outside the model's own text. /p/other/roadmap (empty project) said Nothing is planned. Building now was not exercised (no runs in the clone). Loading the page once crashed the daemon via /api/milestones' concurrent map write, and renders took 25-51s.

## Links

### Out

- verifies [web-roadmap-milestone-page](../capability/rhiza/web-monitor/web-roadmap-milestone-page.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
