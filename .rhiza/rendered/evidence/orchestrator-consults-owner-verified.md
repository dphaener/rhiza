# orchestrator-consults-owner-verified

evidence · `01M785QPDABJ30Y996394YZHZ0` · actor [claude-code-c8d2334f](../actor/claude-code-c8d2334f.md)

## result

PASS

## statement

At 3652345eb ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestOrchestratorConsultsOwner|TestOrchestratorDocRules' -v: both PASS. Read the diff: agents/orchestrator.md 'Consult the owner first' now says consult with a recommendation before any judgment call and act autonomously only when the owner has granted autonomy for that task or session; a new 'When idle' section says propose the next high-value work instead of starting it while queued bug requests still flow through plan-bug/bugfix. grep of AGENTS.md finds no consult/recommendation rule (only the bug-workflow autonomy text). The test checks each phrase and the absence of 'Consult the owner' in AGENTS.md. Note for review: the new 'When idle' heading now sits above the existing 'Quote the owner exactly' and 'Plan a milestone only with the owner' paragraphs, so they read as part of the idle section.

## Links

### Out

- verifies [orchestrator-consults-owner](../capability/rhiza/integrity/own-agent-docs/orchestrator-consults-owner.md)

## History

- `01M76NX8F3G4KXGHQ9G93TEH4Z` agents/orchestrator.md: how the orchestrator runs Rhiza's own development, consulting the owner before judgment calls and acting autonomously only when granted
