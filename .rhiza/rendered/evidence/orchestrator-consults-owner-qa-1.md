# orchestrator-consults-owner-qa-1

evidence · `01M785YVGV1F8TNVP7R48H5CBN` · actor [claude-code-eaa7dd4e](../actor/claude-code-eaa7dd4e.md)

## result

PASS

## statement

QA at d34deebf1: read agents/orchestrator.md: 'Consult the owner first' requires a recommendation before any judgment call and autonomy only when the owner granted it for that task or session; 'When idle' says propose the next high-value work instead of starting it while queued bug requests still flow through plan-bug and bugfix. grep finds no 'Consult the owner' in AGENTS.md. go test -trimpath -count=1 ./cmd/rhiza -run TestOrchestratorConsultsOwner passes; rhiza check on a clean detached checkout of the head prints ok.

## Links

### Out

- verifies [orchestrator-consults-owner](../capability/rhiza/integrity/own-agent-docs/orchestrator-consults-owner.md)

## History

- `01M76NX8F3G4KXGHQ9G93TEH4Z` agents/orchestrator.md: how the orchestrator runs Rhiza's own development, consulting the owner before judgment calls and acting autonomously only when granted
