# run-start-refuses-agent-session-qa-1

evidence · `01M65G0S0VAN50RBVWPRVN7PY2` · actor [claude-code-46d3a9a6](../actor/claude-code-46d3a9a6.md)

## result

PASS

## statement

At 87e2c0a I built rhiza with go build outside PATH and, in a scratch git repo with a change set and a one-step command workflow (RHIZA_ROOT/CHANGE/RUN unset), ran run start --workflow one inside this agent session: exit 2, 'run start refused: this is agent session 46d3a9a6-9056-4c25-beb7-01a18aff4b33, and no agent session starts a run' with next 'ask the owner to run it outside an agent session'; the same with --json and with api call run.start (a usage error naming the session), and with --dry-run. No .git/rhiza was created and git worktree list showed no run worktree. With RHIZA_SESSION unset the same start ran step a, succeeded, wrote .git/rhiza/runs/<id>.jsonl and the cs/run-starts worktree; a later start inside the session was still refused with no new journal. A whitespace-only RHIZA_SESSION counts as no session. go test ./cmd/rhiza -run TestScripts/run-start-agent-session passed.

## Links

### Out

- verifies [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
