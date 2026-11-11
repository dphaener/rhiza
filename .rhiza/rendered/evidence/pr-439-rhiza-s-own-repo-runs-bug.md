# pr-439-rhiza-s-own-repo-runs-bug

evidence · `01M75WC3472J4HDFD3NK2CM8FW` · actor [dphaener](../actor/dphaener.md)

## head

5e8a3c1c1f3e24b0ac0060be13340e86d4ad3110

## pr

https://github.com/dphaener/rhiza/pull/439

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #439 (https://github.com/dphaener/rhiza/pull/439) passed CI at 5e8a3c1

## status

verified

## Links

### Out

- verifies [own-repo-bug-workflows](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-workflows.md)
- verifies [own-repo-bug-flow-note](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-flow-note.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
