# run-spawn-conduct-qa-2

evidence · `01M7BVNWW14V5XH8V2QHZV62R7` · actor [claude-code-348d450a](../actor/claude-code-348d450a.md)

## result

PASS

## statement

Scratch build: outside a run rhiza run spawn exits 2 naming RHIZA_RUN; in a live run a missing --brief exits 2 naming the path; spawned check children ran in the run's worktree under step conduct with their own brief.md and slice.json, and the spawn printed the child's outcome JSON (outcome pass, role, check, step, parent).

## Links

### Out

- verifies [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
