# agent-step-checks-qa-2

evidence · `01M7DTKPF3E5HYKG33Y900M59H` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: go test ./internal/workflow -run TestAgentStepChecks and ./internal/api -run 'TestRunFindingsAgentStepChecks|TestRunFindingsCheckNames' pass; the built binary's workflow describe ship-change shows conduct declaring roles plan, code, qa; run findings --check plan outside a run is refused naming RHIZA_RUN (exit 2).

## Links

### Out

- verifies [agent-step-checks](../capability/rhiza/workflows/built-in-workflows/built-in-review/agent-step-checks.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
