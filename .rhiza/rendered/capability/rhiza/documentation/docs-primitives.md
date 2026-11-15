# docs-primitives

capability · `01M7E5H4GA0EHHSE7MZ91QF875` · parent [documentation](../documentation.md)

## acceptance

A test in cmd/rhiza (docs_primitives_test.go, new) fails unless docs/primitives.md has a section for each of the eight primitives and names every step kind workflow.Kinds (or the step-kind list rhiza docs renders) reports; go test -trimpath ./cmd/rhiza -run 'TestDocsPrimitives|TestScripts/(docs-commands-exist|docs-links)' passes; make self-check passes (rhiza check finds every rhiza command the page names and every link resolves). Touches docs/primitives.md (new), README.md, docs/runs-and-daemon.md, docs/custom-workflows.md and cmd/rhiza/docs_primitives_test.go.

## statement

docs/primitives.md is the primitives reference in the docs' plain style: for each of agent, spawn, gate, push/pr, wait, merge, signoff and hook it says what it does, the step kind or rhiza command that is it (agent step, rhiza run spawn, gate step, rhiza ship open-pr, wait step, rhiza ship merge, signoff step, command step and [[cleanup]]), its outcomes and how a red result hands back to an agent step (handback:<step>, handback_cap); it links docs/reference/workflow-file.md for the keys, and the README, docs/runs-and-daemon.md and docs/custom-workflows.md link it.

## Links

### Out

- decided-by [primitives-doc-hand-written](../../../decision/primitives-doc-hand-written.md)

### In

- depends-on [docs-bring-your-own-agent](docs-bring-your-own-agent.md)
- verifies [docs-primitives-verified](../../../evidence/docs-primitives-verified.md)
- verifies [docs-primitives-qa-1](../../../evidence/docs-primitives-qa-1.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
