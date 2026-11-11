# own-repo-bug-workflows-verified

evidence · `01M75VDHZX56G1G0N1BS5Z9K1Y` · actor [claude-code-f246576c](../actor/claude-code-f246576c.md)

## result

PASS

## statement

At 15e6a4da, with a fresh build (go build -trimpath ./cmd/rhiza) and RHIZA_ROOT/CHANGE/RUN unset: rhiza workflow check plan-bug and check bugfix exit 0 (valid, each with only the expected warning about autonomous runs with bypassed permissions); workflow describe plan-bug shows '(repo)' .rhiza/workflows/plan-bug.toml, includes plan, autonomy plan/accept/ship/signoff auto; describe bugfix shows '(repo)', includes build, review, ship-change, all autonomy auto; workflow ls lists bugfix and plan-bug as repo from .rhiza/workflows. make self-check in a clean detached worktree at 15e6a4da exits 0 (in the run checkout it fails only on render of the verifier's own untracked actor record). Diff read: .rhiza/dispatch.toml sets [request] bug = "plan-bug" and [ready] bug = "bugfix" with no schedules; each workflow is one kind = "workflow" step with autonomous = true; no intake script, eventCommand hook or issue label added.

## Links

### Out

- verifies [own-repo-bug-workflows](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-workflows.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
