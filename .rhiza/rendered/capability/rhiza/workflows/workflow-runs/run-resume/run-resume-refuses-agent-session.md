# run-resume-refuses-agent-session

capability · `01M6AG9P18J4KJEAWFPQ5N731G` · parent [run-resume](../run-resume.md)

## acceptance

A testscript stops a run, sets RHIZA_SESSION and expects rhiza run resume to fail with the session ID on stderr and the run's journal unchanged; with RHIZA_SESSION cleared the resume goes on

## statement

rhiza run resume refuses inside an agent session (RHIZA_SESSION set) exactly as rhiza run start does, naming the session on stderr and journaling nothing

## Links

### Out

- decided-by [resume-refuses-agent-session](../../../../../decision/resume-refuses-agent-session.md)

### In

- constrains [no-agent-session-starts-runs](../../../../../constraint/no-agent-session-starts-runs.md)
- verifies [run-resume-refuses-agent-session-verified](../../../../../evidence/run-resume-refuses-agent-session-verified.md)
- verifies [run-resume-refuses-agent-session-qa-1](../../../../../evidence/run-resume-refuses-agent-session-qa-1.md)
- verifies [pr-159-no-agent-session-starts-runs](../../../../../evidence/pr-159-no-agent-session-starts-runs.md)

## History

- `01M65G0WKX7Y74QCRB8EAN573P` A repo constraint states that no agent session starts a run, so every run-starting path (run start, run resume, a future daemon) is checked against it, not only rhiza run start
