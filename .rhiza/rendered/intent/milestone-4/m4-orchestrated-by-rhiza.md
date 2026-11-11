# m4-orchestrated-by-rhiza

intent · `01M61R4RCY7AE9N46F2YDNX0WD` · parent [milestone-4](../milestone-4.md)

## statement

Milestone 5's first change sets are shipped with rhiza daemon orchestrating and no orchestrator session: the owner accepts plans, answers the inbox and signs off, nothing else.

## Links

### Out

- decided-by [answer-8cd68bc7-6](../../decision/answer-8cd68bc7-6.md)

### In

- serves change set `01M75RDBYXWSMQM0YPEQQ4ECY3` (Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.)

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
