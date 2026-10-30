# ask-status-is-question-verified-2

evidence · `01M66FVF3C4QPWRZPMTHPW2H45` · actor [claude-code-afe5c2a4](../actor/claude-code-afe5c2a4.md)

## result

PASS

## statement

At 04a1df9, go test ./cmd/rhiza -run 'TestScripts/run-blocked$' -count=1 passes: the run's AskUserQuestion, now followed by Claude Code's permission_prompt Notification, reports exactly one ]7501;state=blocked:kind=question with msg decoding to '... Blocked runs: Which colour?', and no kind=permission on stderr; run status shows '- Question: Which colour?'. With the pre-fix internal/api/session_run.go restored in a scratch worktree the same script fails (no '- Question: Which colour?'), so the new questionPending guard is what keeps the question. Diff read: recordHook no longer lets a blocked/permission report replace a pending blocked/question. make fmt-check vet test pass; make self-check fails only on stale .rhiza/rendered actor files from run writes, not on this leaf's code.

## Links

### Out

- verifies [ask-status-is-question](../capability/ask-status-is-question.md)

## History

- `01M61GW8H09EPJ1YCR7P2X3E8J` Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'.
