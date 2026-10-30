# rebase-conflict-usage-limit-pause

question · `01M65PAVNBHXC3KDH2WW77B6P1`

## statement

Rebasing branch cs/usage-limit-pause onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: internal/api/session_ask_test.go, internal/api/session_interview_test.go?

## Links

### In

- answers [usage-limit-rebase-resolution](../decision/usage-limit-rebase-resolution.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
