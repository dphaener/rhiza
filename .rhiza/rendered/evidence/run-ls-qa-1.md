# run-ls-qa-1

evidence · `01M6993C33KTQTHJJGDP5F82QK` · actor [claude-code-84d1030f](../actor/claude-code-84d1030f.md)

## result

PASS

## statement

QA at 0f0c797 with a scratch build on this repo's real runs: rhiza run ls --limit 5 printed five runs newest started first (two running live runs at build/build and review/review, then failed and succeeded ones with no step), each with ID, state, workflow, start time, change set and title; --json gave id, change, title, state, current, started_at, live, workflow, autonomous; --limit 0 gave the default 20; --limit 1001 and -1 exited 2 with 'out of range: give 1 to 1000'; rhiza api call run.ls {"limit":1} returned the newest run only.

## Links

### Out

- verifies [run-ls](../capability/rhiza/workflows/workflow-runs/run-ls.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
