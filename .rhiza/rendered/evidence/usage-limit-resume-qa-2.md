# usage-limit-resume-qa-2

evidence · `01M65N3GDF1RAJGNXXQ20H8E35` · actor [claude-code-a2bc176a](../actor/claude-code-a2bc176a.md)

## result

PASS

## statement

QA cycle 2, built rhiza binary in a scratch repo with a fake claude on PATH: (1) agent reports 'usage limit reached|<past time>' on its first launch: session run relaunched it as `claude --resume <same id>` with the resume prompt, the second launch finished, exactly one start and one resume, exit code 0. (2) agent limited on every launch: 24 resumes, then the session is recorded as error 'claude-code is still stopped at its usage limit after 24 pauses' instead of exited 0. (3) agent left open at its prompt after the limit line: killed at the 30-second transcript poll and resumed. (4) `rhiza session hook Stop` with a transcript ending in 'hit your limit' records working 'paused at its usage limit'; without it, idle. Separately, a working agent whose last line is a tool_use grepping 'hit your limit' was also killed and paused (reported as a finding).

## Links

### Out

- verifies [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
