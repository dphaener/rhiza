# docs-custom-workflow-verified-2

evidence · `01M6ZJSGTBQFHT7P1ZSPRXWEXF` · actor [claude-code-aaaa2ffa](../actor/claude-code-aaaa2ffa.md)

## result

PASS

## statement

Verifier at 5d26135b: go test -trimpath ./cmd/rhiza -run TestScripts/custom-workflows-doc$ passes; it extracts docs/custom-workflows.md's bash blocks, writes the bugfix workflow (reproduce/fix agents, test command with max_visits and on.fail loop, ci wait step) and dispatch.toml [ready] bug = "bugfix", and rhiza workflow check bugfix reports valid. go test ./internal/doclinks passes. The doc walks through workflow new/check/describe and mapping bug requests, and names the generated reference docs/reference/workflow-file.md in plain text plus rhiza guide workflows, which decision custom-workflow-doc-reference-plain-path accepts for this change set. Spot-checked keys (until/timeout/interval/max_visits, --from, dispatch tables ready/request) against internal/workflow and internal/dispatch. Prose reads human; owner read remains with the owner.

## Links

### Out

- verifies [docs-custom-workflow](../capability/rhiza/documentation/docs-custom-workflow.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
