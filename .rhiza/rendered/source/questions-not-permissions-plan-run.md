# questions-not-permissions-plan-run

source · `01M6505PF0CR3E1B9V6DPBPQSE`

## quotes

```json
[]
```

## raw_sha256

e93a5bca7a00e86171178ef2e2352c4108a2e6920f38d8c3f1b585e1aec2ce3c

## session

b0ce8a13-0f72-4b9b-934d-01a18a012c19

## summary

Plan run on change set 01M61GW8H (Questions not permissions), autonomous, so no owner interview. The intent: Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'. Code read: the Claude Code adapter maps AskUserQuestion to blocked with kind question (internal/adapter/claudecode/claudecode.go:162); run status and the run's OSC 7501 report carry kind question (internal/api/run.go:1641, internal/api/session_run.go:502); the run-blocked script asserts state=blocked:kind=question (cmd/rhiza/testdata/script/run-blocked.txtar:61). The Rex program status spec lists question as a valid kind (superlogical docs, build/program-status). The 'needs your permission' wording is not produced in this repo, so it is parked as a question for the owner.

## title

Questions not permissions: plan run, no owner present

## Links

### In

- derived-from [ask-status-is-question](../capability/ask-status-is-question.md)

## History

- `01M61GW8H09EPJ1YCR7P2X3E8J` Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'.
