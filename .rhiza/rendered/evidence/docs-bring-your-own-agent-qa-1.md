# docs-bring-your-own-agent-qa-1

evidence · `01M7EAC21B394A39N2A9XD4RZ6` · actor [claude-code-0cae456c](../actor/claude-code-0cae456c.md)

## result

PASS

## statement

QA with a fresh scratch build: rhiza workflow check passes on docs/examples/byo/skill.toml, conductor.toml and workflow.toml (exit 0, only not-installed skill warnings); a copy with a bad step kind fails (exit 1). Followed the page in a scratch project: workflow new myteam-ship --from ship-change-classic and myteam-conduct --from ship-change work, the forked conduct step matches the page's shape, skill.toml copied in checks by name, change workflow myteam-ship pins and an unknown name is refused. skill check help has --role, --keep, --agent and a required --permission-mode. The page is linked from README.md, docs/agents.md and docs/custom-workflows.md. TestDocsExampleWorkflows and TestDocsBringYourOwn pass; make self-check passes in a clean detached worktree.

## Links

### Out

- verifies [docs-bring-your-own-agent](../capability/rhiza/documentation/docs-bring-your-own-agent.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
