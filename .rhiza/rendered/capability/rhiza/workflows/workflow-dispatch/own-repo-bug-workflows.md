# own-repo-bug-workflows

capability · `01M75RK16E77QY68C4GTVWBE19` · parent [workflow-dispatch](../workflow-dispatch.md)

## acceptance

rhiza workflow check plan-bug and rhiza workflow check bugfix pass (a warning about autonomous runs with bypassed permissions is expected), rhiza workflow describe plan-bug and bugfix each show the repo layer and autonomous, rhiza workflow ls lists both from .rhiza/workflows, and make self-check passes

## statement

This repo dispatches bug requests to autonomous workflows: .rhiza/dispatch.toml sets [request] bug = "plan-bug" and [ready] bug = "bugfix" with no schedules, .rhiza/workflows/plan-bug.toml wraps the built-in plan and .rhiza/workflows/bugfix.toml wraps the built-in ship-change, each as one workflow step with autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; there is no GitHub intake script, no rhiza.eventCommand hook and no issue label

## Links

### Out

- decided-by [bug-workflows-serve-m4](../../../../decision/bug-workflows-serve-m4.md)
- derived-from [bug-workflows-plan](../../../../source/bug-workflows-plan.md)

### In

- depends-on [own-repo-bug-flow-note](own-repo-bug-flow-note.md)
- verifies [own-repo-bug-workflows-verified](../../../../evidence/own-repo-bug-workflows-verified.md)
- verifies [own-repo-bug-workflows-qa-1](../../../../evidence/own-repo-bug-workflows-qa-1.md)
- verifies [pr-439-rhiza-s-own-repo-runs-bug](../../../../evidence/pr-439-rhiza-s-own-repo-runs-bug.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
