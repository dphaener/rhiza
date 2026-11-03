# headless-run-agent-tabs-qa-2

evidence · `01M6GNDAJ3G74M8EQE8B6DFPKD` · actor [claude-code-b014619b](../actor/claude-code-b014619b.md)

## result

PASS

## statement

At fdbebe1e I built rhiza with go build into a scratch dir and ran it by hand in a scratch git repo, with RHIZA_* and REX_* unset and a fake rex that logs its calls (no agent ever answers). Every rex window new in rex.log came from an agent session: 'window new -s rhiza-qa qa: Built · run <run> · build/build · builder built-leaf #1' for each ship run, and 'window new -s rhiza-qa qa: Add a widget · run <run> · interview' for the plan run. Each named the repo session, the run and the step, and rex.log had no split at all.

## Links

### Out

- verifies [headless-run-agent-tabs](../capability/rhiza/runtime/rex-sessions/headless-run-agent-tabs.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
