# usage-limit-resume-qa-1

evidence · `01M65K3JZBZFS1YBRMN63DJFVE` · actor [claude-code-f6599a51](../actor/claude-code-f6599a51.md)

## result

PASS

## statement

QA at 88b6eae with the built binary (go build ./cmd/rhiza), in a scratch repo with a scripted fake claude on PATH and a fake rex: rhiza session start --prompt 'do the step', then rhiza session run <id>. (1) Fake writes 'Claude AI usage limit reached|1700000000' (reset in the past) on its first launch: session run relaunched at once as 'claude --resume <same id>' with prompt 'The usage limit has reset. Continue where you left off.', exactly one start and one resume, ended 'exited with code 0'. (2) Reset 6 s ahead: the resume came 6 s after the start, i.e. it waited for the given reset. (3) 'You've hit your limit' with no reset time: session run was still waiting after 8 s (15-minute default). (4) Limit reported on every launch: 1 start + 24 resumes, then ended 'exited with code 0', status exited code 0.

## Links

### Out

- verifies [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
