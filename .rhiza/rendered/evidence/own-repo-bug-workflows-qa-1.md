# own-repo-bug-workflows-qa-1

evidence · `01M75W09FQDDFPJB3ASC0MDFQW` · actor [claude-code-70b49af7](../actor/claude-code-70b49af7.md)

## result

PASS

## statement

QA with a binary built from the branch head (7ef51341): rhiza workflow check plan-bug and check bugfix each print valid with exit 0 and only the expected bypassed-permissions warning; rhiza workflow describe plan-bug shows (repo) .rhiza/workflows/plan-bug.toml, includes plan, autonomy plan/accept/ship/signoff auto; describe bugfix shows (repo) bugfix.toml, includes build, review, ship-change, all autonomy auto; rhiza workflow ls lists bugfix and plan-bug as repo from .rhiza/workflows; .rhiza/dispatch.toml has [request] bug = plan-bug and [ready] bug = bugfix with no schedules; make self-check in a clean worktree at HEAD prints ok, exit 0

## Links

### Out

- verifies [own-repo-bug-workflows](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-workflows.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
