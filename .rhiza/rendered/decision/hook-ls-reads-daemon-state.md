# hook-ls-reads-daemon-state

decision · `01M6BZV8T9EGJ3QE3ZG02AXKY3`

## statement

The daemon records its notify and event commands and where each came from (daemon flag or git config) in its state at start; rhiza hook ls lists notify and event with command and source, preferring a live daemon's values over git config, and unset otherwise; rhiza hook run <notify|event> [--command c] runs the command once with a sample payload through the daemon's own runner, prints its output and exit status, fails when it fails, exits 2 when unset, and writes nothing to the daemon's state

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m18-orchestrator-picks](../source/m18-orchestrator-picks.md)

### In

- decided-by [hook-ls](../capability/rhiza/hooks/hook-commands/hook-ls.md)
- decided-by [hook-run](../capability/rhiza/hooks/hook-commands/hook-run.md)

## History

- `01M6BVT9HNN801JPDD99CM049F` Plan milestone 18: test your hooks, discuss a change, picked by the orchestrator under the owner's standing instruction
