# init-asks-agent

capability · `01M7E5H4BW9J22TDPF2EBN0XBA` · parent [repo-init](../repo-init.md)

## acceptance

go test -trimpath ./internal/api -run 'TestInit' covers --agent writing rhiza.project.agent, an unknown --agent refused with the registered names and no .rhiza created, the prompt driven through a fake terminal (default on an empty line, a re-ask on an unknown name, a chosen name stored) and no prompt when stdin is not a terminal; go test -trimpath ./cmd/rhiza -run 'TestScripts/(init|root-init|describe-golden)' passes, init.txtar extended for --agent, the non-terminal default and the no-git case, describe-golden -update'd. Touches internal/api/init.go, internal/api/init_test.go, cmd/rhiza (only if the terminal check must live in the CLI frame), cmd/rhiza/testdata/script/init.txtar, describe-golden.txtar and docs/reference/cli.md (rhiza docs); names no agent outside internal/adapter.

## statement

rhiza init asks which agent the project uses: --agent <adapter> answers it (an unknown name is a usage error listing the registered adapters, and nothing is created); with no flag, not --json, and a terminal on stdin it prints the registered adapters with the default marked and reads one line (empty takes the default, an unknown name asks again); with no terminal it takes the default without asking; inside a git repository it writes the answer to rhiza.project.agent, outside one it writes no setting and says how to set it later with rhiza daemon set; its output and --json response name the agent.

## Links

### Out

- depends-on [project-agent-setting](../../runtime/agent-adapters/project-agent-setting.md)
- decided-by [init-asks-agent-setting](../../../../decision/init-asks-agent-setting.md)

### In

- depends-on [docs-bring-your-own-agent](../../documentation/docs-bring-your-own-agent.md)
- verifies [init-asks-agent-verified](../../../../evidence/init-asks-agent-verified.md)
- verifies [init-asks-agent-qa-1](../../../../evidence/init-asks-agent-qa-1.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
