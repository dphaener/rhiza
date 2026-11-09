# docs-custom-workflow-qa-3

evidence · `01M6ZKGP0ME4H6KW62J67DD6MC` · actor [claude-code-71363a3e](../actor/claude-code-71363a3e.md)

## result

PASS

## statement

At 3b5eb86c, with the built binary in a scratch project: rhiza workflow new bugfix, then the doc's cat > .rhiza/workflows/bugfix.toml heredoc; rhiza workflow check bugfix says valid (exit 0), workflow describe bugfix shows reproduce -> fix -> test (fail -> fix) -> ci as described, and run start --workflow bugfix --dry-run with a current change set prints the run. The [ready] bug = "bugfix" dispatch.toml is accepted. The reference is named as the plain path docs/reference/workflow-file.md, per decision custom-workflow-doc-reference-plain-path. TestScripts/custom-workflows-doc and internal/doclinks pass.

## Links

### Out

- verifies [docs-custom-workflow](../capability/rhiza/documentation/docs-custom-workflow.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
