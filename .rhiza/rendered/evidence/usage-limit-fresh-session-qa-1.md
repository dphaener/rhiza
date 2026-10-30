# usage-limit-fresh-session-qa-1

evidence · `01M65K3KP7DM4Q2K027V1FCX1T` · actor [claude-code-f6599a51](../actor/claude-code-f6599a51.md)

## result

PASS

## statement

QA at 88b6eae with the built binary in a scratch repo, scripted fake claude: first launch writes 'Claude AI usage limit reached|<now+4s>' and deletes its own transcript 2 s later, during the pause. rhiza session run waited ~4 s, then relaunched 'claude --session-id <same id>' (not --resume) with the first prompt 'do the step', and printed '- Warning: the transcript of session <id> is gone after the usage-limit pause: started a fresh session for the step instead of resuming'. The fresh start shows only in that printed warning: rhiza session show <id> and the session record hold nothing about it.

## Links

### Out

- verifies [usage-limit-fresh-session](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-fresh-session.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
