# headless-run-agent-tabs-qa-1

evidence · `01M6G7QQ4KAVDS5B5D847K0BFX` · actor [claude-code-5fcf5cce](../actor/claude-code-5fcf5cce.md)

## result

PASS

## statement

At a5ee2b13 I built rhiza (go build) and ran it by hand in a scratch git repo with RHIZA_* and REX_* unset, a fake rex that logs its calls and runs the agent command, and a fake claude agent (held by a file). In the daemon-started headless runs, each agent session (builder, check alpha, fixer) called 'rex window new -s rhiza-repo' titled e.g. 'repo: Built · run <run> · build/build · builder built-leaf #1' and 'repo: One · run <run> · review · fixer #2'. No rex split call appeared. After the held agents were released, rex.log showed 'rex block close -s rhiza-repo block:b5/b6/b7' for exactly the three blocks those window new calls returned, and the journal recorded pane-closed per session. A --foreground run with REX_SESSION=session:mine opened its tab in session:mine (the caller's session) and closed only that block.

## Links

### Out

- verifies [headless-run-agent-tabs](../capability/rhiza/runtime/rex-sessions/headless-run-agent-tabs.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
