# project-agent-setting

capability · `01M7E5H49QWATQ109YQMJKD8EX` · parent [agent-adapters](../agent-adapters.md)

## acceptance

go test -trimpath ./internal/api -run 'TestProjectAgent' covers the setting unset (the default adapter), set (run start's agent steps, session start and skill check take it), an explicit --agent overriding it, and daemon set refusing an unknown adapter name with the registered names in the message; go test -trimpath ./cmd/rhiza -run 'TestScripts/(project-agent|docs-reference-config|describe-golden)' passes, project-agent.txtar new, after rhiza docs regenerates docs/reference/config.md and describe-golden is -update'd if a request type changed. Touches internal/api/daemon_settings.go, internal/api/session.go, internal/api/run.go (run start's agent), internal/api/skill_check.go, docs/reference/config.md and cmd/rhiza/testdata/script/project-agent.txtar; names no agent outside internal/adapter.

## statement

A project setting rhiza.project.agent (git config, listed by rhiza daemon settings and in docs/reference/config.md) names the agent adapter the project uses; rhiza daemon set refuses a name the adapter registry does not know, listing the registered names; rhiza run start, rhiza session start and rhiza skill check use it when --agent is unset, an explicit --agent still wins, and unset keeps the registered default adapter.

## Links

### Out

- decided-by [init-asks-agent-setting](../../../../decision/init-asks-agent-setting.md)

### In

- depends-on [init-asks-agent](../../model/repo-init/init-asks-agent.md)
- verifies [project-agent-setting-verified](../../../../evidence/project-agent-setting-verified.md)
- verifies [project-agent-setting-qa-1](../../../../evidence/project-agent-setting-qa-1.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
