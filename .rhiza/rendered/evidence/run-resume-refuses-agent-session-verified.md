# run-resume-refuses-agent-session-verified

evidence · `01M6AX47D7HMHSTSFSHP9VQAEA` · actor [claude-code-178fd615](../actor/claude-code-178fd615.md)

## result

PASS

## statement

At 23d866e5 ran go test ./cmd/rhiza -run TestScripts/run-resume-agent-session -count=1 -v: PASS. The script stops a run (run start --no-wait), sets RHIZA_SESSION=01SESSIONFORTESTING, and rhiza run resume exits 2 with stderr 'run resume refused: this is agent session 01SESSIONFORTESTING, and no agent session starts a run'; cmp of .git/rhiza/runs/* before/after shows the journal unchanged; with RHIZA_SESSION cleared the resume continues and the run succeeds. Diff read: refuseAgentSession in internal/api/run.go uses the same message and next hint as run start's refusal, and sits in runRunResume at the same point start puts it (after the dry-run return, before canRun and any journal write). go test ./..., go vet and fmt-check pass; make check's only failure is a stale render for this verifier's own actor record, not part of the commit.

## Links

### Out

- verifies [run-resume-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-resume/run-resume-refuses-agent-session.md)

## History

- `01M65G0WKX7Y74QCRB8EAN573P` A repo constraint states that no agent session starts a run, so every run-starting path (run start, run resume, a future daemon) is checked against it, not only rhiza run start
