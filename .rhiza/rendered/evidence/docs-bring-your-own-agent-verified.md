# docs-bring-your-own-agent-verified

evidence · `01M7E8T8G5VZFXN80EB0MCY2AE` · actor [claude-code-524c84b7](../actor/claude-code-524c84b7.md)

## result

PASS

## statement

Verified at e323b6cbe in a detached checkout. A scratch build's rhiza workflow check passed on docs/examples/byo/skill.toml, conductor.toml and workflow.toml (exit 0 each; the only warnings are that the /myteam: skills are not installed). go test -p 2 -trimpath ./cmd/rhiza -run 'TestDocsExampleWorkflows|TestDocsBringYourOwn|TestScripts/docs-commands-exist' passed, and so did the Checks: variant with docs-links (no such script exists, so that pattern matches nothing). make self-check exited 0. TestDocsExampleWorkflows now stats and checks all three byo files; TestDocsBringYourOwn requires the page to link each example and primitives.md, and requires README.md, docs/agents.md and docs/custom-workflows.md to link the page. I checked the page's claims against the code. rhiza skill check defaults to the builder role and takes --role verifier|check; --permission-mode is required and bypass modes are refused unless the machine is trusted as dedicated; --keep exists; it exits 1 when a gate fails; its scratch dir comes from os.MkdirTemp. run spawn always prompts 'follow rhiza guide <role>' (run_spawn.go:178). {brief}, {slice} and {handback} are brief.md, slice.json (the rhiza slice --json bytes) and handback.md (run_brief.go). The four ship-change gate commands and the plan, code and qa checks match builtin/ship-change.toml. A declared check with no round fails review-check (conductorReviewCheck). ship-change-classic includes the build workflow, whose leaves step has builder and verifier tables. invoke is an alias of prompt. rhiza init --agent and rhiza.project.agent work as described: an explicit --agent wins, unset falls back to the registered default, the claude-code adapter is the only one, and daemon set's validateAgent lists the known names. The custom-workflows#point-bugs-at-it anchor exists. The diff touches only the leaf's listed files plus its .rhiza records, and leaves no TODOs.

## Links

### Out

- verifies [docs-bring-your-own-agent](../capability/rhiza/documentation/docs-bring-your-own-agent.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
