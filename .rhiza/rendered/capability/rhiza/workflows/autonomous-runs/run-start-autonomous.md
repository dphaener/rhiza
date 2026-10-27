# run-start-autonomous

capability · `01M5WRHDKXZ0Z9PQKMPRVQBC5W` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

A testscript starts a run with --autonomous on a fake adapter: the journal's started entry has autonomous true, run status --json and its markdown show autonomous, the fake agent's environment has RHIZA_AUTONOMOUS=1, and after a stop run resume (no flag) launches the next agent with it too; a run started without the flag has none of these. describe-golden is regenerated for the new request and response fields.

## statement

rhiza run start --autonomous starts a run in autonomous mode: the journal's started entry records it, run resume keeps it without the flag, run status shows it, and every agent the run launches has RHIZA_AUTONOMOUS=1 in its environment.

## Links

### Out

- decided-by [autonomous-mode](../../../../decision/autonomous-mode.md)
- decided-by [autonomous-on-session-record](../../../../decision/autonomous-on-session-record.md)

### In

- depends-on [autonomous-keeps-permissions](autonomous-keeps-permissions.md)
- depends-on [autonomous-gate-rule](autonomous-gate-rule.md)
- depends-on [arbiter-on-question](../arbiter-sessions/arbiter-on-question.md)
- verifies [run-start-autonomous-verified](../../../../evidence/run-start-autonomous-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
