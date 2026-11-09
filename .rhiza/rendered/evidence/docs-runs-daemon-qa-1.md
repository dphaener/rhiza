# docs-runs-daemon-qa-1

evidence · `01M6ZJE698JKGYKNACNMGV7BQ8` · actor [claude-code-80c03c5f](../actor/claude-code-80c03c5f.md)

## result

FAIL

## statement

At d23439f6, with the built binary in a scratch repo: every rhiza command in docs/runs-and-daemon.md exists (--help); guide install, rex install --force, run send-back --note, request --kind/--ref match; daemon settings lists every named setting; daemon set rhiza.autonomy.{plan,accept,ship,signoff} owner works; the doc's [[schedules]] weekly-deps block in dispatch.toml is accepted (rhiza health: '1 schedule ... next run 2026-11-09T09:00:00Z'). But: daemon settings shows defaults autonomy.accept/ship/signoff = auto, contradicting the doc's 'By default the daemon asks you at every gate'; rhiza request and run start fail 'is not inside a Rhiza root / next: rhiza root init' until rhiza root init + project add, which the doc's setup never mentions; and workflow-file autonomy (autonomous = true|false or [autonomy] table, internal/workflow/workflow.go parseAutonomy) is not described, though the leaf names per-workflow autonomy.

## Links

### Out

- verifies [docs-runs-daemon](../capability/rhiza/documentation/docs-runs-daemon.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
