# init-asks-agent-verified

evidence · `01M7E89GDXN1GKXRPKPXXWNBGQ` · actor [claude-code-23ee3b45](../actor/claude-code-23ee3b45.md)

## result

PASS

## statement

In a detached worktree at f12face57: go test -trimpath ./internal/api -run TestInit passed (TestInit, TestInitAgent: --agent written to rhiza.project.agent, unknown --agent exit 2 listing the registered adapters with no .rhiza and no setting, no terminal takes the default, prompt through a fake terminal: empty line takes the default, unknown name re-asked, chosen name stored, --json asks nothing, no git writes nothing); go test -trimpath ./cmd/rhiza -run 'TestScripts/(init|root-init|describe-golden|project-agent|docs-.*)$' passed without -update (golden and cli.md current); go vet ./cmd/rhiza ./internal/api clean. A scratch build driven through a real pty (script -qec) listed claude-code (default), refused 'nope' and asked again, took the default on an empty line and wrote rhiza.project.agent=claude-code; with stdin /dev/null outside git it asked nothing, wrote no setting and printed the rhiza daemon set hint; init --agent bad exited 2 with the adapter list and created nothing. Diff: the terminal check is in cmd/rhiza/ctx.go (Ctx.Terminal in internal/api/api.go), no agent name in non-test code outside internal/adapter, no TODOs; project-agent.txtar adjusted for the setting init now writes.

## Links

### Out

- verifies [init-asks-agent](../capability/rhiza/model/repo-init/init-asks-agent.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
