# bug-workflows-serve-m4

decision · `01M75RK10JHEXK5AHEDTK8TMVR`

## statement

Rhiza's own bug workflows serve m4-orchestrated-by-rhiza: they take the owner out of planning and shipping a bug, as recent fixes to Rhiza's own runs serve milestone 4's intents

## wording

(assumed by the planner in an autonomous run, no owner to ask)

## Links

### Out

- derived-from [bug-workflows-plan](../source/bug-workflows-plan.md)

### In

- decided-by [own-repo-bug-workflows](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-workflows.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
