# project-agent-setting-qa-1

evidence · `01M7EAC1TG42NV1N7Z6ZR3RA4H` · actor [claude-code-0cae456c](../actor/claude-code-0cae456c.md)

## result

PASS

## statement

QA with a fresh scratch build at e323b6cbe in a hermetic scratch repo: rhiza daemon set rhiza.project.agent bogus exits 2 listing 'want one of claude-code' and writes nothing; claude-code is written and daemon settings lists it (git config, project), unset lists it as "" (default, project). With the git config set to bogus behind daemon set's back, run start --dry-run (step build), session start and skill check each refuse naming the key; an explicit --agent claude-code gets past it (run start would start); set and unset both start. go test -p 2 -trimpath ./internal/api -run 'TestInit|TestProjectAgent|TestSession' and TestScripts/(project-agent|docs-reference-config|describe-golden) pass.

## Links

### Out

- verifies [project-agent-setting](../capability/rhiza/runtime/agent-adapters/project-agent-setting.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
