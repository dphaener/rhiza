# bug-workflows-plan

source · `01M75RJJ1ZKJH3SB2FKBAHHWFV`

## quotes

```json
[]
```

## raw_sha256

ce08d45635388bf80bdc5c74aa02c2d983b4971d6e5ef9c265c5dcd47b970855

## session

ad605dc2-cc4e-4522-b7a2-01a1cb8746da

## summary

Autonomous plan run (no owner present) for change set 01M75RDBYXWSMQM0YPEQQ4ECY3. The change set's intent already carries the owner's request and scope, quoted there from 2026-11-11: copy Nebari's bug dispatch (nebari PR #94) into this repo, minus the GitHub issue intake. The planner read Nebari's .rhiza/dispatch.toml, plan-bug.toml and bugfix.toml, this repo's .rhiza/workflows (docs.toml, ship.sh; no ship-change override, no dispatch.toml), docs/runs-and-daemon.md and docs/custom-workflows.md (which already document dispatch in general). Assumptions taken without escalation: the change set serves m4-orchestrated-by-rhiza, as recent fixes to Rhiza's own runs serve milestone 4's intents; the config is one leaf (dispatch.toml and both workflows together, since dispatch naming missing workflows is broken); the docs note goes in CLAUDE.md, the repo's contributor and agent guidance, because the product docs already explain dispatch and this note is about this repo only. No questions parked.

## title

Autonomous plan of Rhiza's own bug workflows

## Links

### In

- derived-from [bug-workflows-serve-m4](../decision/bug-workflows-serve-m4.md)
- derived-from [bug-flow-note-in-claude-md](../decision/bug-flow-note-in-claude-md.md)
- derived-from [own-repo-bug-workflows](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-workflows.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
