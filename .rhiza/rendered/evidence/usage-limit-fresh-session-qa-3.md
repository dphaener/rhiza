# usage-limit-fresh-session-qa-3

evidence · `01M65P8ZWBMYCHQM98P5XZBWNG` · actor [claude-code-193e157c](../actor/claude-code-193e157c.md)

## result

PASS

## statement

QA cycle 3, hands-on with a scratch build: rhiza run start --workflow one (a single agent step) in a scratch repo, with a fake rex that runs the tab command and a fake claude. On first launch the fake wrote the synthetic usage-limit line (reset 4s ahead) and deleted the transcript 2s later. After the wait, the run relaunched claude with --session-id, not --resume: a fresh session for the step. The step passed. The run journal (.git/rhiza/runs/<run>.jsonl) gained a fresh-session entry for step build and that session, and the session record's fresh list got a timestamp. session run printed the warning that the transcript was gone and a fresh session was started instead of resuming.

## Links

### Out

- verifies [usage-limit-fresh-session](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-fresh-session.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
