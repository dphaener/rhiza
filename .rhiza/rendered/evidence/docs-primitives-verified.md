# docs-primitives-verified

evidence · `01M7E77ZGWBMFFW34DF54K00VW` · actor [claude-code-d2f4dc9e](../actor/claude-code-d2f4dc9e.md)

## result

PASS

## statement

In a detached checkout at d7e32101c: go test -p 2 -trimpath ./cmd/rhiza -run 'TestDocsPrimitives|TestScripts/(docs-commands-exist|docs-links)' exit 0 (TestDocsPrimitives and TestScripts/docs-commands-exist pass; no docs-links script exists, so that alternative matches nothing, and links are checked by rhiza check instead); make self-check exit 0. TestDocsPrimitives checks a ## section per primitive (Agent, Spawn, Gate, Push/PR, Wait, Merge, Sign-off, Hook), every workflow.Kinds entry named in backticks, handback:, handback_cap, rhiza run spawn, rhiza ship open-pr, rhiza ship merge, [[cleanup]], the workflow-file.md link, and links from README.md, docs/runs-and-daemon.md and docs/custom-workflows.md. The diff touches only the files the acceptance names plus the builder's log and rendered actor. Claims checked against code: DefaultHandbackCap 3 and the parked stop (internal/run/engine.go), handback.md contents (run/brief.go RenderHandback), live message vs resume and the required handback text (api/run.go stepRunner.Handback), the uncapped-loop warning (api/workflow_check.go), handback targets must be agent steps (workflow.go checkTargets), kind outcomes (workflow.go kindOutcomes), gate runs every command and keeps each failing output in a file (engine.go gate), StopExit 75, max_interval doubling on unchanged output, run spawn flags, prompt, maxRuns/holds and exit 1 on fail (api/run_spawn.go), ship open-pr body replacement (api/ship_openpr.go), signoff on.changes required, and ship-change.toml's steps and routes. Minor note, not a defect: the Agent section's 'passes when the agent says done or exits 0' repeats the code's own agent kind summary (workflow/spec.go), though the engine passes a step without until only on exit; and nothing forbids an agent step's own on.fail = handback:<other agent step>, so 'doesn't hand back itself' describes ship-change practice rather than a rule.

## Links

### Out

- verifies [docs-primitives](../capability/rhiza/documentation/docs-primitives.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
