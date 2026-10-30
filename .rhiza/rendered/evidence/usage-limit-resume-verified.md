# usage-limit-resume-verified

evidence · `01M65JJZCV1VSMH1SV2QYJKFWK` · actor [claude-code-cfad59fb](../actor/claude-code-cfad59fb.md)

## result

PASS

## statement

At commit 6c0cd2b ran go test ./cmd/rhiza -run TestScripts/session-run-usage-limit: a fake claude writes 'Claude AI usage limit reached|1700000000' as the last assistant line on its first launch and 'done' on the second. rhiza session run ended with the agent exited with code 0, and claude.log held exactly one 'start --session-id 158fe87f-…' and one 'resume --resume 158fe87f-…' line: the same stored session ID resumed with the Claude CLI's --resume flag, and no second session ID was started. TestCommandResume (argv uses --resume, never --session-id) and TestUsageLimit pass. Read the diff: session_run.go loops on the same rec.ID, with a.Command(launch) given Resume=true, at most 24 pauses. make check passed (exit 0) on a rerun; the first run failed in the test step, which a plain go test ./... did not reproduce.

## Links

### Out

- verifies [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
