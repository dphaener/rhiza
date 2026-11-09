# docs-runs-daemon-qa-2

evidence · `01M6ZK15VNC19B6CMVE8DA0E3G` · actor [claude-code-acfcd74a](../actor/claude-code-acfcd74a.md)

## result

FAIL

## statement

At 91524236, built rhiza (go build -trimpath) and in a scratch HOME/root: the cycle-1 fixes hold. rhiza daemon settings shows plan=owner, accept/ship/signoff=auto as the page now says; rhiza request fails 'is not inside a Rhiza root' until rhiza root init + project add, then queues; in a copy of ship-change, autonomous = true / false gives all four stages 'auto/owner from workflow mine' in run start --dry-run, an [autonomy] signoff/ship = owner table overrides only those two (the rest 'from rhiza.autonomy.<stage>'), and both together make workflow check report 'autonomous and [autonomy] are both set'; run request-changes <file> [<run>] matches; the commands named answer --help; doclinks passes. But rhiza run start --workflow ship-change (and --workflow mine --dry-run, the page's own example) fails 'step "build/build" is a leaves step whose builder runs with no permission mode: give --permission-mode', even after rhiza daemon set rhiza.daemon.permissionMode acceptEdits; it gets past the check only with --permission-mode, which the page never mentions. Also project add ~/code/myapp fails 'not inside a Rhiza root' when run from outside the root.

## Links

### Out

- verifies [docs-runs-daemon](../capability/rhiza/documentation/docs-runs-daemon.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
