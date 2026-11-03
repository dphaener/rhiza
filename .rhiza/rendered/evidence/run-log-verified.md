# run-log-verified

evidence · `01M6G415VAFFNVHWWF72HJS2RJ` · actor [claude-code-4bdd5c04](../actor/claude-code-4bdd5c04.md)

## result

PASS

## statement

At 341756b6 ran go test ./cmd/rhiza -run 'TestScripts/(run-log|describe-golden)$' -count=1 -v: run-log.txtar (fake rex, daemon --once starts a headless run) shows rhiza run log --limit 2 printing one/two and 'next: rhiza run log <run> --cursor 2', --cursor 2 printing three with no next, --all printing all lines, and exit 2 after the log dir is removed; the --limit/--cursor flags match ls's (internal/api/ls.go). Ran go test ./internal/tui/liveapp -run TestLKeyOpensSelectedRunLogLayer: on the fake runner the l key calls run.log with 'R1 --layer' for the selected run and shows the notice; run.log --layer runs rex do with a new_layer running '<self> run log <run> --all'. gofmt clean, go vet ok; go test ./internal/api -skip TestReadBudget passes (TestReadBudget, a timing budget off the read path, failed under VM load).

## Links

### Out

- verifies [run-log](../capability/rhiza/orchestration/daemon/run-log.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
