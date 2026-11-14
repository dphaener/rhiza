# web-roadmap-milestone-data-qa-1

evidence · `01M7BN416T6VMR4SRTW5VA33GM` · actor [claude-code-3866e93d](../actor/claude-code-3866e93d.md)

## result

FAIL

## statement

QA at 33cc5a571, scratch build run on a clone of the branch's repo: rhiza milestones --json gave 36 milestones, 6 finished; every done point's met flag and every milestone's state word matched the rules (checked with jq), and each done point's change sets were in its milestone's list. --view bogus and a short what-if ref failed cleanly with exit 2. But 2 of 11 runs died with 'fatal error: concurrent map writes' at internal/api/milestones.go:164 (an unlocked map write in the EachInFlight callback), and the same crash killed a daemon serving /roadmap.

## Links

### Out

- verifies [web-roadmap-milestone-data](../capability/rhiza/web-monitor/web-roadmap-milestone-data.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
