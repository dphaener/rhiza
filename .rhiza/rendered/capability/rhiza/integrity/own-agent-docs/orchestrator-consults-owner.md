# orchestrator-consults-owner

capability · `01M76P0AJ9VF4D377JGWZJAAZY` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

The orchestrator-doc test finds in agents/orchestrator.md the consult-with-a-recommendation rule, autonomy only when the owner granted it, and the idle rule (propose, don't start; bug requests still flow), and finds no consult rule in AGENTS.md; go test -trimpath ./cmd/rhiza -run on its name passes

## statement

agents/orchestrator.md tells the orchestrator to consult the owner, with a recommendation, before any judgment call and to act autonomously only when the owner granted autonomy for that task or session; when idle it proposes the next high-value work to the owner instead of starting it, while queued bug requests still flow through their autonomous workflows; AGENTS.md does not carry this rule

## Links

### In

- verifies [orchestrator-consults-owner-verified](../../../../evidence/orchestrator-consults-owner-verified.md)
- verifies [orchestrator-consults-owner-qa-1](../../../../evidence/orchestrator-consults-owner-qa-1.md)
- verifies [pr-496-orchestrator-agent-doc](../../../../evidence/pr-496-orchestrator-agent-doc.md)

## History

- `01M76NX8F3G4KXGHQ9G93TEH4Z` agents/orchestrator.md: how the orchestrator runs Rhiza's own development, consulting the owner before judgment calls and acting autonomously only when granted
