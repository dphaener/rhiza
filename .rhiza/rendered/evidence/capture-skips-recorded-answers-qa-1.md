# capture-skips-recorded-answers-qa-1

evidence · `01M6MN7ASPTCTYTHRJEV4Q7NET` · actor [claude-code-2253c893](../actor/claude-code-2253c893.md)

## result

FAIL

## statement

QA at 628808db with a scratch build (go build -trimpath) in a hermetic scratch repo (own HOME, fake rex, a 3-question AskUserQuestion transcript). In the change set: Decision question='Which port?' (verbatim), 'Which log format? ' (trailing space), and 'Retry count?' (live). session interview --dry-run and --json proposed only answer-<sid8>-2 (Which log format?) and listed answers 1 and 3 under skipped ('Skipped ...: already recorded'); a rerun proposed nothing. After retiring the 'Retry count?' Decision, rhiza session done on a new session wrote only answer-<sid8>-3, so retired Decisions do not stop capture. Landed path: record change set with source add --raw-session, a Decision 'Which port?' derived-from that Source, then rhiza land; session interview skipped the port answer and proposed only the unrecorded one. But session done's response (decisions, warnings) carries no skipped list (captureAfterRun in session_run.go:458 drops got.Skipped), so session done never reports a skipped answer as the statement and acceptance ask. go test ./cmd/rhiza -run 'TestScripts/session-(interview|done)': ok.

## Links

### Out

- verifies [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
