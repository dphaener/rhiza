# session-token-use

capability · `01M5WRWXYCPD79ST702XJ65GAF` · parent [session-journal](../session-journal.md)

## acceptance

A testscript runs a session in a run with a fake agent that leaves a transcript fixture where TranscriptGlobs points, says done, adds to the transcript and exits: session show --json and its markdown give the counts of the later reading, and the run's journal has a tokens entry at done and one at exit with the step, session and counts. A Go engine test shows each agent step's end reports its session once, including a session left open. With no transcript, session done and session run warn, session run exits with the agent's code, and neither the record nor the journal gains token use. describe-golden is regenerated for the new response field.

## statement

When a session says done, when its step ends, and again when its agent exits, rhiza reads its token use from the transcript, records it on the session record and, when the session belongs to a run, appends a tokens entry with the step, session and the four counts of the session and of its helper agents to the run's journal; the later reading wins. rhiza session show shows it. A transcript that cannot be read is a warning and records no token use.

## Links

### Out

- depends-on [transcript-token-use](transcript-token-use.md)
- decided-by [m3-scope](../../../../decision/m3-scope.md)
- decided-by [helper-token-total](../../../../decision/helper-token-total.md)
- decided-by [token-use-read-at-done-and-step-end](../../../../decision/token-use-read-at-done-and-step-end.md)

### In

- depends-on [run-status-tokens](run-status-tokens.md)
- verifies [session-token-use-tests](../../../../evidence/session-token-use-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
