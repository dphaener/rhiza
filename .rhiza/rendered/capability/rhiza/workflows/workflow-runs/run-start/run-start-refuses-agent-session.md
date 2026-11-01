# run-start-refuses-agent-session

capability · `01M6505GVBMV8B71YA23043CEE` · parent [run-start](../run-start.md)

## acceptance

A testscript sets RHIZA_SESSION and runs rhiza run start on a change set: it exits non-zero, the error names the session ID, and no run journal or worktree is written; without RHIZA_SESSION the start is unchanged

## statement

rhiza run start refuses inside an agent session, when RHIZA_SESSION is set, and the error names that session

## Links

### Out

- decided-by [refuse-run-start-in-agent-session](../../../../../decision/refuse-run-start-in-agent-session.md)

### In

- constrains [no-agent-session-starts-runs](../../../../../constraint/no-agent-session-starts-runs.md)
- blocks [dry-run-in-agent-session](../../../../../question/dry-run-in-agent-session.md)
- blocks [resume-in-agent-session](../../../../../question/resume-in-agent-session.md)
- verifies [run-start-refuses-agent-session-verified](../../../../../evidence/run-start-refuses-agent-session-verified.md)
- verifies [run-start-refuses-agent-session-qa-1](../../../../../evidence/run-start-refuses-agent-session-qa-1.md)
- verifies [run-start-refuses-agent-session-verified-2](../../../../../evidence/run-start-refuses-agent-session-verified-2.md)
- verifies [pr-104-agents-never-start-runs](../../../../../evidence/pr-104-agents-never-start-runs.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
