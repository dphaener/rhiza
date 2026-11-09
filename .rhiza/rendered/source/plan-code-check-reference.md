# plan-code-check-reference

source · `01M6ZTD50JQ9GWCB71EVRG3CSY`

## quotes

```json
[]
```

## raw_sha256

ba83644b627efe93231c8a2018bc46ef09832751cd5da84c6b91cfd44f717176

## session

e5dcdc72-7dbf-4e17-9520-01a1bfa42f36

## summary

Autonomous plan run for change set 01M6ZT7 (bug: review/code fails every run that edits CLI help, settings or the workflow format). The planner read internal/api/run_codecheck.go (checkIn runs runCheck in-process) and internal/api/check.go (DocsWant from docsFiles() and Commands from Operations(), both the running binary's), and found the docs-command rule has the same flaw as the reference staleness check: both judge the branch with the installed binary. One question went to the arbiter: which intent to serve, and how to fix it (skip both rules in the code check; build the worktree's rhiza; or skip only when the build differs). The arbiter decided (code-check-leaves-reference-to-self-check): serve m30-reference, since #412 introduced the regression; the code check turns off the generated-reference comparison and the docs-command rule through an explicit CheckRequest option, plain rhiza check and make self-check keep both, and the code check does not build the worktree's binary. The planner wrote one leaf for that behaviour.

## title

Plan: code check leaves the generated reference to self-check

## Links

### In

- derived-from [code-check-leaves-reference-to-self-check](../decision/code-check-leaves-reference-to-self-check.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
