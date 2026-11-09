# question-e5dcdc72-1

question · `01M6ZTAKPAXQ86RT3PX7CQZECJ`

## session

e5dcdc72-7dbf-4e17-9520-01a1bfa42f36

## statement

Change set 01M6ZT7 (review/code check fails every run that edits CLI help, settings or the workflow format, because the in-process rhiza check compares docs/reference with the installed binary's rhiza docs): which intent should it serve? It is a regression from #412, the generated reference. / How should run code-check stop judging the branch's docs/reference (and the docs-command rule, which has the same flaw for a newly added, documented command) with the installed binary instead of the branch's build?

## Links

### In

- answers [code-check-leaves-reference-to-self-check](../decision/code-check-leaves-reference-to-self-check.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
