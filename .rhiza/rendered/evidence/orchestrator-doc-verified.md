# orchestrator-doc-verified

evidence · `01M785N1J9R4R38BFVRG3D26GZ` · actor [claude-code-8b6f4051](../actor/claude-code-8b6f4051.md)

## result

PASS

## statement

At b60dc8980, go test -trimpath ./cmd/rhiza -run 'TestOrchestratorDocRules|TestAgentDocs' -v -count=1 passes: TestOrchestratorDocRules (new, cmd/rhiza/orchestratordoc_test.go) reads agents/orchestrator.md and finds each rule: no git fetch in the main clone, no pkill by pattern, nothing uncommitted on main / every model write on its own branch and PR, rhiza.daemon.maxRuns under load, merge origin/main and never rebase+force-push, gh pr merge --squash, quote the owner exactly, plan a milestone only with the owner; TestAgentDocsConsistent and its fixtures pass with orchestrator.md indexed in agents/README.md. Read the diff: the doc states each rule with its reason; gofmt and go vet clean.

## Links

### Out

- verifies [orchestrator-doc](../capability/rhiza/integrity/own-agent-docs/orchestrator-doc.md)

## History

- `01M76NX8F3G4KXGHQ9G93TEH4Z` agents/orchestrator.md: how the orchestrator runs Rhiza's own development, consulting the owner before judgment calls and acting autonomously only when granted
