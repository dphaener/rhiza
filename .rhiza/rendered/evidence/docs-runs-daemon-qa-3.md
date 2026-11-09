# docs-runs-daemon-qa-3

evidence · `01M6ZKGNYJWA6W8C7AMPDXWJ8Q` · actor [claude-code-71363a3e](../actor/claude-code-71363a3e.md)

## result

PASS

## statement

At 3b5eb86c, with the built binary in a scratch root: rhiza root init <root>, cd <root>, rhiza project add <root>/myapp registers the project (from outside the root it still fails, as the doc now warns). run start <change> --workflow ship-change --dry-run without --permission-mode stops asking for one, even with rhiza.daemon.permissionMode set; with --permission-mode acceptEdits it prints the run; bypassPermissions is refused without trust = "dedicated" in machine.toml, as the doc says. daemon settings shows plan=owner and accept/ship/signoff=auto, matching the Autonomy section. A workflow copied with workflow new mine --from ship-change plus an [autonomy] table (signoff, ship = owner) passes workflow check, and run start --workflow mine --permission-mode acceptEdits --dry-run prints ship and signoff as 'owner from workflow mine', the rest from rhiza.autonomy.<stage>. run request-changes <file> [<run>] matches the doc. Every rhiza command the page names answers --help; internal/doclinks passes.

## Links

### Out

- verifies [docs-runs-daemon](../capability/rhiza/documentation/docs-runs-daemon.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
