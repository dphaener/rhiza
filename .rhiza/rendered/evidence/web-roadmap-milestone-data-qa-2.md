# web-roadmap-milestone-data-qa-2

evidence · `01M7BVHJ0QFFWM9XHXG9DXCBJP` · actor [claude-code-e789a95e](../actor/claude-code-e789a95e.md)

## result

PASS

## statement

QA at c1769e866, a fresh scratch build run on a clone of the branch: rhiza milestones --json ran 12 times, every one exit 0 with byte-identical output (36 milestones, 6 finished), no 'concurrent map writes' crash and none in a daemon serving the pages; every done point's met flag and every milestone's state word matched the rules (checked with jq). rhiza api call milestones printed the same bytes. The human text listed milestones, done points [x]/[ ] and '6 of 36 finished'. --view as-is and what-if:<this change set> worked; --view bogus, an unknown what-if ref and an extra argument failed cleanly with exit 2; an empty project gave 'No milestones yet.' and finished 0.

## Links

### Out

- verifies [web-roadmap-milestone-data](../capability/rhiza/web-monitor/web-roadmap-milestone-data.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
