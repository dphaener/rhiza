# orchestrator-doc-qa-1

evidence · `01M785YVDX4G163CH0JVG03Q4Q` · actor [claude-code-eaa7dd4e](../actor/claude-code-eaa7dd4e.md)

## result

PASS

## statement

QA at d34deebf1 with a binary built from the head (go build -trimpath): go test -trimpath -count=1 ./cmd/rhiza -run 'TestOrchestratorDocRules|TestOrchestratorConsultsOwner|TestAgentDocs' passes (TestAgentDocsConsistent and Fixtures pass with orchestrator.md indexed in agents/README.md). Exercised the commands the doc names: in a scratch git repo `rhiza daemon set rhiza.daemon.maxRuns 2` exits 0 and writes git config 2, `... abc` exits 2 with a whole-number error; `rhiza withdraw/accept/apply/link --help` all exist; `git ls-remote origin refs/heads/main` works as the fetch-free alternative; the ../AGENTS.md link resolves. In a clean detached worktree of the head, `rhiza check` prints ok (only pre-existing unclaimed warnings) and `rhiza render` leaves no diff.

## Links

### Out

- verifies [orchestrator-doc](../capability/rhiza/integrity/own-agent-docs/orchestrator-doc.md)

## History

- `01M76NX8F3G4KXGHQ9G93TEH4Z` agents/orchestrator.md: how the orchestrator runs Rhiza's own development, consulting the owner before judgment calls and acting autonomously only when granted
