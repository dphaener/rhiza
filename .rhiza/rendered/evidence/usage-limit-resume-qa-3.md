# usage-limit-resume-qa-3

evidence · `01M65P8Z3QRHW7SGS8PAXPB1W6` · actor [claude-code-193e157c](../actor/claude-code-193e157c.md)

## result

PASS

## statement

QA cycle 3, hands-on with a scratch build of the branch head (rhiza session run in a scratch repo, fake claude and rex on PATH, isolated HOME). The agent wrote Claude Code's synthetic usage-limit line with a reset 6s ahead. The run waited 6s, and session show read working, 'paused at its usage limit until <reset>', during the wait. It then relaunched claude with --resume and the same session ID, told 'The usage limit has reset. Continue where you left off.' as the last arg. No second --session-id launch, exit 0. An agent left open at its limit (sleeping) was killed by the 30s watcher and resumed (elapsed 30s). An agent limited every time got 1 start and 24 resumes, then the session failed: 'still stopped at its usage limit after 24 pauses'. Via rhiza session hook: a Stop hook with the synthetic 'hit your limit' line last recorded working 'paused at its usage limit'. A normal Stop, and a later model turn quoting 'usage limit reached|1700000000' and 'hit your limit', recorded idle.

## Links

### Out

- verifies [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
