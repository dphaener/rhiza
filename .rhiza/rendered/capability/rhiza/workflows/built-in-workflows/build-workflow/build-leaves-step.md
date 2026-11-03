# build-leaves-step

capability · `01M5YX3BA0JVJ2PHZ2N4TX9SX9` · parent [build-workflow](../build-workflow.md)

## acceptance

Parse tests load a leaves step with both sub-tables and reject one missing builder or verifier, an unknown sub-table key, `models` on a step kind other than agent or leaves, `split` on a step kind other than leaves, and a bad `split` value; an engine test with a fake runner over two leaves finds a journal entry per leaf session, attempt and outcome, and `rhiza run status` shows them per leaf.

## statement

Workflow files take a `leaves` step kind with `builder` and `verifier` sub-tables (prompt, agent, permission_mode, models) and a `split` setting, and the engine runs it as one step that conducts the run change set's leaves, journaling each leaf's sessions, attempts and outcome.

## Links

### Out

- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [engine-is-conductor](../../../../../decision/engine-is-conductor.md)
- decided-by [answer-ab54d544-1](../../../../../decision/answer-ab54d544-1.md)
- decided-by [build-ruling-1](../../../../../decision/build-ruling-1.md)
- decided-by [build-ruling-7](../../../../../decision/build-ruling-7.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-leaf-order](build-leaf-order.md)
- depends-on [build-leaf-commit](build-leaf-commit.md)
- depends-on [build-verifier-findings](build-verifier-findings.md)
- depends-on build-panes (retired)
- depends-on [run-program-status](../../workflow-runs/run-program-status.md)
- verifies [build-leaves-step-verified](../../../../../evidence/build-leaves-step-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
