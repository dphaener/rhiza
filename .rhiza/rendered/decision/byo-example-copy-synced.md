# byo-example-copy-synced

decision · `01M7RB61WHX6VTXBDJSF4VT2YE`

## statement

The BYO daemon scripts run a copy of docs/examples/byo/workflow.toml (myteam-flow) from their archive, picked by .rhiza/dispatch.toml [ready], held byte-identical to the docs file by a Go test; docker, make, gh, rex and claude are fakes on PATH and /myteam:implement is a stub skill under $HOME/.claude/skills

## wording

(planner, from the code: a testscript cannot read repo files; dispatch [ready] is how the daemon picks a non-default workflow)

## Links

### Out

- derived-from [byo-daemon-plan-research](../source/byo-daemon-plan-research.md)

### In

- decided-by [byo-workflow-daemon-e2e](../capability/rhiza/workflows/workflow-runs/byo-workflow-daemon-e2e.md)
- decided-by [byo-spawn-hold-intervention](../capability/rhiza/workflows/workflow-runs/byo-spawn-hold-intervention.md)
- decided-by [byo-failed-run-retry](../capability/rhiza/workflows/workflow-runs/byo-failed-run-retry.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
