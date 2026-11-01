# run-resume-refuses-agent-session-qa-1

evidence · `01M6AXJQ11D3E8SKFE8J8XZ6K8` · actor [claude-code-88ab6ba0](../actor/claude-code-88ab6ba0.md)

## result

PASS

## statement

Hands-on at e8198352 with a scratch build (go build -o $TMPDIR/rhiza-x ./cmd/rhiza) in a fresh scratch git repo, RHIZA_* env unset: rhiza init, change new, a one-step wait workflow, run start --no-wait -> run stopped. With RHIZA_SESSION set to a real session ID, rhiza run resume <run> exits 2 with stderr 'run resume refused: this is agent session <id>, and no agent session starts a run' and next 'ask the owner to run it outside an agent session'; --json and rhiza api call run.resume return the same usage error with the ID. shasum of .git/rhiza/runs/* before and after all refused attempts is identical (journal unchanged). --dry-run inside the session still previews (as run start --dry-run does); a whitespace-only RHIZA_SESSION is not treated as a session. With RHIZA_SESSION unset the same resume continues and the run passes. run start in the session still refuses with the same wording. make check: fmt, vet and go test pass; self-check fails only on stale renders for this review run's actor records.

## Links

### Out

- verifies [run-resume-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-resume/run-resume-refuses-agent-session.md)

## History

- `01M65G0WKX7Y74QCRB8EAN573P` A repo constraint states that no agent session starts a run, so every run-starting path (run start, run resume, a future daemon) is checked against it, not only rhiza run start
