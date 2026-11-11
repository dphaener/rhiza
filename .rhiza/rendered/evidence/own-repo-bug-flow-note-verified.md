# own-repo-bug-flow-note-verified

evidence · `01M75VJ3DYWBK5B9AVXE9SNN27` · actor [claude-code-ecaf6bbd](../actor/claude-code-ecaf6bbd.md)

## result

PASS

## statement

Verifier claude-code (opus) at da28284e: read the diff; CLAUDE.md gains a '## Bug requests' section naming `rhiza request --kind bug`, the CLI/Request app/orchestrator/watchdog, `.rhiza/dispatch.toml`, plan-bug and bugfix, says both are autonomous with no owner step whatever rhiza.autonomy.* says, and that bugs never come from GitHub issues. Matches .rhiza/dispatch.toml ([request] bug=plan-bug, [ready] bug=bugfix) and `autonomous = true` in both workflow files. With a scratch build of this commit (RHIZA_ROOT/CHANGE/RUN unset): `rhiza request --help` exit 0 and lists --kind feature|bug|milestone; `rhiza workflow describe plan-bug` and `rhiza workflow describe bugfix` exit 0 from this checkout's .rhiza/workflows, all autonomy auto; `rhiza workflow check plan-bug`/`bugfix` are valid with the bypassed-permissions warnings the note calls expected.

## Links

### Out

- verifies [own-repo-bug-flow-note](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-flow-note.md)

## History

- `01M75RDBYXWSMQM0YPEQQ4ECY3` Rhiza's own repo runs bug requests through autonomous bug workflows, as Nebari does (nebari PR #94: .rhiza/dispatch.toml, .rhiza/workflows/plan-bug.toml and bugfix.toml in ~/code/nebari), without Nebari's GitHub issue intake. Owner, 2026-11-11: "i want to add a bugfix workflow like nebari has as well" and "it won't come from github issues though". Bugs reach this repo as rhiza request --kind bug (CLI, the Request app, the orchestrator, the watchdog). In this repo: (1) .rhiza/dispatch.toml: [request] bug = plan-bug and [ready] bug = bugfix; no schedules; (2) .rhiza/workflows/plan-bug.toml wraps the built-in plan and bugfix.toml wraps the built-in ship-change (Rhiza has no ship-change override), both autonomous = true, so a bug is planned and shipped with no owner step whatever rhiza.autonomy.* says; (3) a short docs note on how bug requests flow. Out of scope: no github-intake script, no rhiza.eventCommand hook, no issue labels. rhiza workflow check is expected to warn about autonomous workflows with bypassed permissions. Test with rhiza workflow check and make self-check.
