# usage-limit-fresh-session-qa-2

evidence · `01M65N3H4BX491FK1HY46V93PA` · actor [claude-code-a2bc176a](../actor/claude-code-a2bc176a.md)

## result

PASS

## statement

QA cycle 2, built rhiza binary with a fake claude: the agent reports the limit with a reset 5s ahead, and its transcript is removed during the pause. session run started a fresh `claude --session-id <same id>` (not --resume), the step ended with code 0, the output warned 'the transcript ... is gone after the usage-limit pause: started a fresh session for the step instead of resuming', and the session record in .git/rhiza/sessions kept a 'fresh' timestamp. `rhiza session show` (human and --json) does not display that timestamp.

## Links

### Out

- verifies [usage-limit-fresh-session](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-fresh-session.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
