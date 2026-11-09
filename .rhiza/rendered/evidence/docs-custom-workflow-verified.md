# docs-custom-workflow-verified

evidence · `01M6ZJ6A5GKSCJYXVYN5T2PP1Z` · actor [claude-code-acf11ee9](../actor/claude-code-acf11ee9.md)

## result

PASS

## statement

At 040f5f6c, verified docs/custom-workflows.md against the leaf's acceptance. With GOTMPDIR outside the repo, `go test -trimpath ./cmd/rhiza -run TestScripts/custom-workflows-doc -count=1` passes. That testscript runs the doc's bash blocks, which write .rhiza/workflows/bugfix.toml and .rhiza/dispatch.toml (bug = "bugfix"). Then `rhiza workflow check bugfix` reports 'workflow bugfix is valid' and describe shows the fix/test loop. `go test -trimpath ./internal/doclinks/... -count=1` passes. The doc walks through writing a bugfix workflow with reproduce, fix, test (bounded by max_visits) and ci steps, maps bug requests to it through dispatch.toml, and names the generated workflow-file reference as the plain path docs/reference/workflow-file.md, alongside `rhiza guide workflows`. Decision custom-workflow-doc-reference-plain-path accepts the plain path for this change set, and Generated reference adds the link once it writes the file. The prose reads naturally, in a direct second-person voice with no stock LLM phrasing, which meets doc-prose-humanized. The owner's read is still to come and is outside a verifier's reach.

## Links

### Out

- verifies [docs-custom-workflow](../capability/rhiza/documentation/docs-custom-workflow.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
