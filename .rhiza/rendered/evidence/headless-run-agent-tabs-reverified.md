# headless-run-agent-tabs-reverified

evidence · `01M6GD3X5HG0TB5G98TT8KZ61C` · actor [claude-code-45ff71bd](../actor/claude-code-45ff71bd.md)

## result

PASS

## statement

At b0c22e78 (review cycle 2: model-only commit retiring the tiling leaves). Ran go test ./cmd/rhiza -count=1 -run 'TestScripts/(headless-run-agent-tabs|run-agent-tabs|build-tabs|arbiter-agent-tabs|daemon-headless-runs)$': ok. headless-run-agent-tabs.txtar has a fake rex and starts the run with rhiza daemon --once as a detached child with no REX_SESSION. It asserts 4 'window new -s rhiza-repo repo: Built · run <id> · build · <role> <leaf> #1' calls, no 'split' in the rex log, and close calls that equal the opened blocks one for one. A second run queued from REX_SESSION=session:mine still opens its tabs in rhiza-repo and never touches block:mine. In the change set view, rhiza show --change confirms run-sessions-one-window, agent-steps-tile-run-window, question-sessions-tile-run-window, checks-and-fixes-anchored, anchor-gone-splits-run-window and build-panes are retired, and the leaf statement names them. make fmt-check vet self-check: exit 0. make test fails only on wall-clock budget tests (internal/api TestReadBudget, internal/slice TestBudget, internal/tui/modelapp TestGeneratedModel; medians 60-90ms against 50ms) under load average 15 on 4 cores with other rhiza runs active. Those tests read generated models that this commit does not touch, and build/check passed at 08:06.

## Links

### Out

- verifies [headless-run-agent-tabs](../capability/rhiza/runtime/rex-sessions/headless-run-agent-tabs.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
