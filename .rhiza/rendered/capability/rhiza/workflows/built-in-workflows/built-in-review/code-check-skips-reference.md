# code-check-skips-reference

capability · `01M6ZTDC04JYQN7GT0XNB3MKFV` · parent [built-in-review](../built-in-review.md)

## acceptance

A testscript in a repo whose rhiza-written docs/reference is stale and whose docs name a command the registry lacks sees rhiza check fail on both, and rhiza run code-check on a run there pass with neither among its causes; a test of the check op shows the option off keeps both rules (go test ./internal/api/... and the touched testscripts; describe-golden regenerated)

## statement

rhiza run code-check runs rhiza check with the generated-reference comparison (docs/reference against rhiza docs) and the docs-command rule turned off, through an explicit CheckRequest option, so a run whose change edits CLI help, settings or the workflow format, or adds a documented command, is not failed by the installed binary's reference; plain rhiza check and make self-check keep both rules

## Links

### Out

- depends-on [review-code-check](review-code-check.md)
- decided-by [code-check-leaves-reference-to-self-check](../../../../../decision/code-check-leaves-reference-to-self-check.md)

### In

- verifies [code-check-skips-reference-verified](../../../../../evidence/code-check-skips-reference-verified.md)
- verifies [code-check-skips-reference-qa-1](../../../../../evidence/code-check-skips-reference-qa-1.md)
- verifies [code-check-skips-reference-verified-3](../../../../../evidence/code-check-skips-reference-verified-3.md)
- verifies [pr-417-fix-every-run-whose-change](../../../../../evidence/pr-417-fix-every-run-whose-change.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
