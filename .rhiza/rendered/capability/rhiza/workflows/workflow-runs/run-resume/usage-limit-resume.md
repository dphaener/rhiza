# usage-limit-resume

capability · `01M6507G6PPPPEVQYKYJNQRSD6` · parent [run-resume](../run-resume.md)

## acceptance

With a fake agent that reports the usage limit once, the run resumes the stored session ID with the Claude CLI's resume flag, the step completes, and no second session starts

## statement

A run whose agent session hits the Claude usage limit continues that same session once the limit resets, instead of failing

## Links

### Out

- decided-by [usage-limit-resume-same-session](../../../../../decision/usage-limit-resume-same-session.md)
- decided-by [usage-limit-detect-and-wait-answer](../../../../../decision/usage-limit-detect-and-wait-answer.md)

### In

- blocks [usage-limit-detect-and-wait](../../../../../question/usage-limit-detect-and-wait.md)
- verifies [usage-limit-resume-verified](../../../../../evidence/usage-limit-resume-verified.md)
- verifies [usage-limit-resume-qa-1](../../../../../evidence/usage-limit-resume-qa-1.md)
- verifies [usage-limit-resume-verified-2](../../../../../evidence/usage-limit-resume-verified-2.md)
- verifies [usage-limit-resume-qa-2](../../../../../evidence/usage-limit-resume-qa-2.md)
- verifies [usage-limit-resume-verified-3](../../../../../evidence/usage-limit-resume-verified-3.md)
- verifies [usage-limit-resume-qa-3](../../../../../evidence/usage-limit-resume-qa-3.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
