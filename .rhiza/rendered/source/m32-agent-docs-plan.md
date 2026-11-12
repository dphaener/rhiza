# m32-agent-docs-plan

source · `01M76NVMK5983DAD0RYY0B0BVG`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I want another feature request. claude code will read agents.md now, so let's make it agents.md. I also want a directory with agent docs, with some indexed instructions for agents. For instance, an orchestrator. You've been doing a great job, but often times you move a little fast and make decisions without consulting with me. I know this can be confusing at times, since I do sometimes ask for full autonomy, but that should only happen when asked for. This is purely for developers working on Rhiza and the agents that work on them. If this seems too big for a single feature, let's make it another milestone.",
    "where": "2026-11-11 request, relayed in the record's intent"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Move to AGENTS.md: AGENTS.md becomes the one source; CLAUDE.md shrinks to a single @AGENTS.md import line",
    "where": "plan interview, CLAUDE.md"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Orchestrator doc only: only agents/orchestrator.md carries the consult rule",
    "where": "plan interview, consult rule"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "A Go test fails when an agents/*.md file is missing from the index or a relative link in AGENTS.md or agents/ is broken",
    "where": "plan interview, checks"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Developing Rhiza, Run agents, Operating the VM",
    "where": "plan interview, agents/ docs"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Prune them: delete each memory entry the docs now cover; personal-only facts stay",
    "where": "plan interview, memory"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Three change sets; all planned now; the orchestrator prunes its memory after land",
    "where": "plan interview, split"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Propose, don't start: when idle, the orchestrator picks the next high-value work and asks the owner with a recommendation; it starts work on its own only when the owner granted autonomy for that task or session. Queued bug requests still flow",
    "where": "plan interview, idle rule"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "None of these: Nebari hands-off, the owner's Rex habits and priorities stay in personal memory",
    "where": "plan interview, scope"
  }
]
```

## raw_sha256

c0c76243f7652b7e0aa1e7861f61939b987bdf58d161230e2b675967f3bb000a

## session

9eeeaa0d-e0b8-40d8-b95a-01a1cce67858

## summary

Planned milestone 32 from the owner's request for AGENTS.md and an indexed agents/ directory of instructions for agents developing Rhiza. All picks were the planner's AskUserQuestion options. Decided: CLAUDE.md's content moves to AGENTS.md and CLAUDE.md shrinks to an @AGENTS.md import; the consult-before-deciding rule lives in agents/orchestrator.md only, not in AGENTS.md; a Go test keeps the agents/README.md index and relative links consistent; agents/ starts with orchestrator.md, developing.md, run-agents.md and vm.md; the orchestrator's personal memory entries those docs cover are pruned, by the orchestrator after each docs change set lands (not a built leaf); three change sets (AGENTS.md entry point; orchestrator doc; the other three role docs, both depending on the first), all planned now; when idle the orchestrator proposes the next high-value work with a recommendation instead of starting it, unless the owner granted autonomy, while bug requests still flow; Nebari hands-off, the owner's Rex habits and priorities stay in personal memory, not in the docs. Nothing left open.

## title

Milestone 32 plan: agent docs for Rhiza's own developers, with the owner

## Links

### In

- derived-from [milestone-32](../intent/milestone-32.01M76NWA4W870AXQ6PBANDWRGE.md)
- derived-from [answer-9eeeaa0d-1](../decision/answer-9eeeaa0d-1.md)
- derived-from [answer-9eeeaa0d-2](../decision/answer-9eeeaa0d-2.md)
- derived-from [answer-9eeeaa0d-3](../decision/answer-9eeeaa0d-3.md)
- derived-from [answer-9eeeaa0d-4](../decision/answer-9eeeaa0d-4.md)
- derived-from [answer-9eeeaa0d-5](../decision/answer-9eeeaa0d-5.md)
- derived-from [answer-9eeeaa0d-6](../decision/answer-9eeeaa0d-6.md)
- derived-from [answer-9eeeaa0d-7](../decision/answer-9eeeaa0d-7.md)
- derived-from [answer-9eeeaa0d-8](../decision/answer-9eeeaa0d-8.md)
- derived-from [answer-9eeeaa0d-9](../decision/answer-9eeeaa0d-9.md)
- derived-from [answer-9eeeaa0d-10](../decision/answer-9eeeaa0d-10.md)

## History

- `01M76ECB9S2Q9QPWRAX3502XGY` Agent docs for Rhiza's own developers. Owner, 2026-11-11: "I want another feature request. claude code will read agents.md now, so let's make it agents.md. I also want a directory with agent docs, with some indexed instructions for agents. For instance, an orchestrator. You've been doing a great job, but often times you move a little fast and make decisions without consulting with me. I know this can be confusing at times, since I do sometimes ask for full autonomy, but that should only happen when asked for. This is purely for developers working on Rhiza and the agents that work on them. If this seems too big for a single feature, let's make it another milestone." The owner then chose (orchestrator's options): a milestone, and the directory at agents/ in the repo root (index at agents/README.md, linked from AGENTS.md). Starting points (verify): CLAUDE.md (112 lines) holds today's repo guidance; docs/agents.md is the user-facing product guide to Rhiza's agents and stays as it is; the product's own guide install into AGENTS.md is a separate area (Area guide), not this milestone. The orchestrator role today lives in the orchestrator's personal memory on the VM (watch loop rules: no git fetch in the main clone, no pkill by pattern, main tree stays clean, nothing uncommitted on main, maxRuns under load, merge-main-never-rebase, squash merges, quoting the owner exactly, milestone planning needs the owner, consult before deciding unless autonomy was granted for the task).
