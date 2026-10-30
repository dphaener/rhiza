# ask-status-is-question-qa-1

evidence · `01M66F8EKNC18P48TV2Z19WSRN` · actor [claude-code-d290bfde](../actor/claude-code-d290bfde.md)

## result

PASS

## statement

QA at e5649b3 with the built binary (go build ./cmd/rhiza), RHIZA_* env unset, in a scratch repo: rhiza run start --workflow ask with a fake claude whose AskUserQuestion ('Which colour?') goes through the built binary's real 'rhiza session hook PreToolUse', and a fake rex for window new. run status read 'blocked' with '- Question: Which colour?' and --json kind question; the run's stderr carried exactly one OSC 7501 ']7501;state=blocked:kind=question:app=rhiza:' report, no kind=permission, and the run succeeded after the answer. That exact report, replayed into a real Rex window, came back from real 'rex events --json' as program_status_changed with state blocked, kind question, msg 'Blocked runs: build waits on you in tab repo: Blocked runs: Which colour?'.

## Links

### Out

- verifies [ask-status-is-question](../capability/ask-status-is-question.md)

## History

- `01M61GW8H09EPJ1YCR7P2X3E8J` Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'.
