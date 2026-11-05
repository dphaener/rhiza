# live-app-runs-qa-3

evidence · `01M69AVCG8H6RGTMBX3D2S8WVZ` · actor [claude-code-b426a4bf](../actor/claude-code-b426a4bf.md)

## result

PASS

## statement

QA at e65e34a with a scratch build in a real 80x24 Rex pane over this repo's 20 runs: rhiza tui app live shows 'Daemon: running · 2 live runs' on line 1 and the key help on line 24, with each run's title, state, step, leaves n/n, checks and elapsed. j x25 scrolled to the last run with the daemon line kept on top; enter there set the notice 'Run ... is not live' and the frame still fit 24 lines. j/k redrew in ~0.1s and q quit in ~0.3s (cycle 2 measured 7-10s). tui app live --once took 0.29s and cmp matches rhiza run ls --json byte for byte. go test -race ./internal/tui/liveapp ./internal/tui/roadmapapp pass.

## Links

### Out

- verifies live-app-runs (retired)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
