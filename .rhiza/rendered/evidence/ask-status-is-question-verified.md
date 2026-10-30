# ask-status-is-question-verified

evidence · `01M66F44DX0WVA9A0D7S2ZR1NJ` · actor [claude-code-8ae0c633](../actor/claude-code-8ae0c633.md)

## result

PASS

## statement

At commit 9dd7934, go test ./cmd/rhiza -run 'TestScripts/run-blocked$' -count=1 passes. The script drives a fake claude whose AskUserQuestion goes through the real PreToolUse hook, and asserts exactly one OSC 7501 report ']7501;state=blocked:kind=question:app=rhiza:' whose msg decodes to 'Blocked runs: build waits on you in tab repo: Blocked runs: Which colour?' (the question at its end), plus '! stderr kind=permission'. The diff adds only those assertions; the code path is claudecode.go mapping AskUserQuestion to blocked/question and run.go writing :kind=question. make check exits 0.

## Links

### Out

- verifies [ask-status-is-question](../capability/ask-status-is-question.md)

## History

- `01M61GW8H09EPJ1YCR7P2X3E8J` Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'.
