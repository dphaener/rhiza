# m32-agents-md

intent · `01M76NWA870PYJG4FH6PF4WHRJ` · parent [milestone-32](../milestone-32.md)

## statement

AGENTS.md at the repo root holds today's CLAUDE.md guidance, CLAUDE.md only imports it, AGENTS.md links agents/README.md, the index of every agents/ doc, and a Go test keeps that index and the docs' relative links consistent.

## Links

### In

- serves change set `01M76NWXV16NF94MBV6WYESRM0` (AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md)

## History

- `01M76ECB9S2Q9QPWRAX3502XGY` Agent docs for Rhiza's own developers. Owner, 2026-11-11: "I want another feature request. claude code will read agents.md now, so let's make it agents.md. I also want a directory with agent docs, with some indexed instructions for agents. For instance, an orchestrator. You've been doing a great job, but often times you move a little fast and make decisions without consulting with me. I know this can be confusing at times, since I do sometimes ask for full autonomy, but that should only happen when asked for. This is purely for developers working on Rhiza and the agents that work on them. If this seems too big for a single feature, let's make it another milestone." The owner then chose (orchestrator's options): a milestone, and the directory at agents/ in the repo root (index at agents/README.md, linked from AGENTS.md). Starting points (verify): CLAUDE.md (112 lines) holds today's repo guidance; docs/agents.md is the user-facing product guide to Rhiza's agents and stays as it is; the product's own guide install into AGENTS.md is a separate area (Area guide), not this milestone. The orchestrator role today lives in the orchestrator's personal memory on the VM (watch loop rules: no git fetch in the main clone, no pkill by pattern, main tree stays clean, nothing uncommitted on main, maxRuns under load, merge-main-never-rebase, squash merges, quoting the owner exactly, milestone planning needs the owner, consult before deciding unless autonomy was granted for the task).
