# run-start-refuses-agent-session-verified-2

evidence · `01M65H5EZFZ7Y7FCFDYH1R2A0T` · actor [claude-code-a688af17](../actor/claude-code-a688af17.md)

## result

PASS

## statement

Attempt 2 at 5c01390. Read the diff: the RHIZA_SESSION guard in runRunStart (internal/api/run.go) now sits after the dry-run return and before canRun, t.open (worktree) and startRun (journal), so nothing is written before the refusal. go test ./cmd/rhiza -run TestScripts/run-start-agent-session passes: it sets RHIZA_SESSION, run start exits non-zero, stderr names the session, no .git/rhiza/runs, and without RHIZA_SESSION the start succeeds. By hand in a scratch repo with a build of this commit: RHIZA_SESSION=01SESSX rhiza run start --workflow one exited 2 with 'run start refused: this is agent session 01SESSX, and no agent session starts a run', no .git/rhiza and no sibling repo.run-starts worktree. --dry-run in the session exited 0 and wrote nothing. With RHIZA_SESSION empty the run succeeded and created .git/rhiza/runs and ../repo.run-starts. go test ./... passed; one TestGeneratedModel timing-budget failure under load passed on a solo rerun. make self-check failed only on stale rendered actor/index files from model writes made during the run, not from this commit.

## Links

### Out

- verifies [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
