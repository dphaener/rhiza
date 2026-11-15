# docs-bring-your-own-agent

capability · `01M7E5H4JPW8SNCH4HE643V6DT` · parent [documentation](../documentation.md)

## acceptance

TestDocsExampleWorkflows (cmd/rhiza/docs_examples_test.go), extended to docs/examples/byo/*.toml, fails unless all three examples (skill.toml, conductor.toml, workflow.toml) pass rhiza workflow check; a test fails unless docs/bring-your-own-agent.md includes or links each example file; go test -trimpath ./cmd/rhiza -run 'TestDocsExampleWorkflows|TestDocsBringYourOwn|TestScripts/docs-commands-exist' passes; make self-check passes. Touches docs/bring-your-own-agent.md (new), docs/examples/byo/ (new), README.md, docs/agents.md, docs/custom-workflows.md and cmd/rhiza/docs_examples_test.go.

## statement

docs/bring-your-own-agent.md, in the docs' plain style, shows a team bringing its own agent at three levels with one worked example each: a skill (rhiza skill check on the team's invoke, then that invoke in a role table of a workflow forked from ship-change-classic), a conductor (the team's conduct skill as ship-change's conduct prompt, reading {brief}, {slice} and {handback}, with the gate's checks it must still pass) and a workflow (the team's own steps from the primitives); it covers rhiza init --agent and rhiza.project.agent, links docs/primitives.md, and is linked from the README, docs/agents.md and docs/custom-workflows.md; each example is a file under docs/examples/byo/.

## Links

### Out

- depends-on [init-asks-agent](../model/repo-init/init-asks-agent.md)
- depends-on [docs-primitives](docs-primitives.md)
- decided-by [byo-levels-skill-conductor-workflow](../../../decision/byo-levels-skill-conductor-workflow.md)

### In

- verifies [docs-bring-your-own-agent-verified](../../../evidence/docs-bring-your-own-agent-verified.md)
- verifies [docs-bring-your-own-agent-qa-1](../../../evidence/docs-bring-your-own-agent-qa-1.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
