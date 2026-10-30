# usage-limit-fresh-session

capability · `01M6507H8DN68SWSSCRT613VZM` · parent [run-resume](../run-resume.md)

## acceptance

With the stored session's transcript removed, the run starts a fresh session for the step and records that it did so

## statement

When a usage-limited session's transcript is missing, the run starts a fresh session for the step instead of resuming

## Links

### Out

- decided-by [usage-limit-resume-same-session](../../../../../decision/usage-limit-resume-same-session.md)

### In

- verifies [usage-limit-fresh-session-verified](../../../../../evidence/usage-limit-fresh-session-verified.md)
- verifies [usage-limit-fresh-session-qa-1](../../../../../evidence/usage-limit-fresh-session-qa-1.md)
- verifies [usage-limit-fresh-session-verified-2](../../../../../evidence/usage-limit-fresh-session-verified-2.md)
- verifies [usage-limit-fresh-session-qa-2](../../../../../evidence/usage-limit-fresh-session-qa-2.md)
- verifies [usage-limit-fresh-session-qa-3](../../../../../evidence/usage-limit-fresh-session-qa-3.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
