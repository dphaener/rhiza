# docs-custom-workflow-qa-1

evidence · `01M6ZJE651XDST69KTXBEZ35BH` · actor [claude-code-80c03c5f](../actor/claude-code-80c03c5f.md)

## result

PASS

## statement

At d23439f6, built rhiza (go build -trimpath) into a scratch dir and, in a fresh git repo under a scratch Rhiza root with its own XDG_CONFIG_HOME: rhiza init; rhiza workflow new bugfix; workflow check bugfix (valid); ran the doc's two bash blocks extracted from docs/custom-workflows.md, which wrote .rhiza/workflows/bugfix.toml and .rhiza/dispatch.toml; rhiza workflow check bugfix printed 'workflow bugfix is valid' exit 0; workflow describe bugfix showed reproduce->fix->test (on.fail -> fix)->ci->done/failed; workflow new mine --from ship-change then check passed; workflow ls listed bugfix and mine as repo; with a current change set, rhiza run start --workflow bugfix --dry-run printed 'Run ... would start' at step reproduce, writing nothing. The reference is named as plain path docs/reference/workflow-file.md per decision custom-workflow-doc-reference-plain-path. go test ./cmd/rhiza -run TestScripts/custom-workflows-doc and ./internal/doclinks pass.

## Links

### Out

- verifies [docs-custom-workflow](../capability/rhiza/documentation/docs-custom-workflow.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
