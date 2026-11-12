# milestone-32

intent · `01M76NWA4W870AXQ6PBANDWRGE`

## milestone

```json
32
```

## statement

Agents developing Rhiza find their instructions in the repo: AGENTS.md is the entry point any agent reads, it links an indexed agents/ directory of role docs (orchestrator, developing Rhiza, run agents, the VM), and those docs replace the rules that lived only in the orchestrator's personal memory.

## Links

### Out

- derived-from [m32-agent-docs-plan](../source/m32-agent-docs-plan.md)

## History

- `01M76ECB9S2Q9QPWRAX3502XGY` Agent docs for Rhiza's own developers. Owner, 2026-11-11: "I want another feature request. claude code will read agents.md now, so let's make it agents.md. I also want a directory with agent docs, with some indexed instructions for agents. For instance, an orchestrator. You've been doing a great job, but often times you move a little fast and make decisions without consulting with me. I know this can be confusing at times, since I do sometimes ask for full autonomy, but that should only happen when asked for. This is purely for developers working on Rhiza and the agents that work on them. If this seems too big for a single feature, let's make it another milestone." The owner then chose (orchestrator's options): a milestone, and the directory at agents/ in the repo root (index at agents/README.md, linked from AGENTS.md). Starting points (verify): CLAUDE.md (112 lines) holds today's repo guidance; docs/agents.md is the user-facing product guide to Rhiza's agents and stays as it is; the product's own guide install into AGENTS.md is a separate area (Area guide), not this milestone. The orchestrator role today lives in the orchestrator's personal memory on the VM (watch loop rules: no git fetch in the main clone, no pkill by pattern, main tree stays clean, nothing uncommitted on main, maxRuns under load, merge-main-never-rebase, squash merges, quoting the owner exactly, milestone planning needs the owner, consult before deciding unless autonomy was granted for the task).
