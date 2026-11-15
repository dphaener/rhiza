# project-agent-setting-verified

evidence · `01M7E6PSVHYXEEEP8EY37S02XC` · actor [claude-code-1ad66a93](../actor/claude-code-1ad66a93.md)

## result

PASS

## statement

In a detached checkout at cf9023b7f: go test -trimpath -p 2 ./internal/api -run 'TestProjectAgent|TestSessionStart|TestSkillCheck|TestDaemonSet|TestWorkflowCheck' passed (TestProjectAgentStart covers unset/default, set, explicit --agent wins and a step's agent wins for session start and run start's agent steps; TestProjectAgentSkillCheck covers skill check unset/set/explicit; TestProjectAgentDaemonSet refuses 'nope' with a usage error listing adapter.Names() and writes nothing, then sets and lists the key; TestProjectAgentUnknown). go test -trimpath ./cmd/rhiza -run 'TestScripts/(project-agent|docs-reference-config|describe-golden|daemon-settings|skill-check)$' passed, project-agent.txtar new. A scratch build's rhiza docs left the tree clean (docs/reference/config.md and cli.md current), gofmt clean. Diff touches the named files plus session_test.go (the refusal test now allows the setting's git config read) and the rendered actor; the added Go code names no agent, using the adapter registry.

## Links

### Out

- verifies [project-agent-setting](../capability/rhiza/runtime/agent-adapters/project-agent-setting.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
