# docs-custom-workflow-qa-2

evidence · `01M6ZK15ZGW9PW4D5C86YB8YZK` · actor [claude-code-acfcd74a](../actor/claude-code-acfcd74a.md)

## result

PASS

## statement

At 91524236 with the built binary: go test -trimpath ./cmd/rhiza -run TestScripts/custom-workflows-doc$ (extracts the doc's bash blocks, writes bugfix.toml and dispatch.toml, rhiza workflow check bugfix valid) and ./internal/doclinks pass; the doc's steps set permission_mode, so its run start --dry-run is not hit by the missing-permission-mode refusal; the generated reference is named as plain path per decision custom-workflow-doc-reference-plain-path.

## Links

### Out

- verifies [docs-custom-workflow](../capability/rhiza/documentation/docs-custom-workflow.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
