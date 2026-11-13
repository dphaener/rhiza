# run-brief-files-qa-2

evidence · `01M7B24D481N743H61MG6Q95JJ` · actor [claude-code-a1cb0102](../actor/claude-code-a1cb0102.md)

## result

PASS

## statement

QA cycle 2: session dirs seen on real runs: build-1/agent, l-1/builder-<leaf>-<attempt>, review-1/<check>-check-1, <check>-critic-1, fixer-<leaf>-<cycle>-<attempt>, fixer-<cycle>-<attempt>, all distinct across cycles and attempts, under .git/rhiza/runs outside the worktree, kept after the run. brief.md has sections in order, both leaves with acceptance, the change set's own decision and constraint (the cycle-2 fix), rhiza.ship.localGate and the contract summary ending 'Full contract: rhiza guide contract'; seven briefs of one model byte-identical. go test ./internal/run -run 'TestBriefFiles|TestSessionDir|TestHandbackFile' ok.

## Links

### Out

- verifies [run-brief-files](../capability/rhiza/workflows/workflow-runs/run-brief-files.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
