# code-check-skips-reference-verified-3

evidence · `01M6ZYV5PC142WNTHZKGHSWEYR` · actor [claude-code-4d4cfa8f](../actor/claude-code-4d4cfa8f.md)

## result

PASS

## statement

At 48896560 (attempt 3): read the diff (CheckRequest.SkipReference with flag skip-reference gates DocsWant/Commands/Docs in runCheck; checkIn passes SkipReference:true; cli.md regenerated). go test -trimpath ./internal/api/... ok; TestCheckSkipReference PASS (default check 2 reference findings, skip_reference 0); TestScripts/review-code-check PASS (rhiza check fails on stale docs/reference/cli.md and `rhiza no-such-command`, run code-check passes with neither among its causes); TestScripts/describe-golden PASS with skip_reference. Branch build's rhiza docs leaves docs/reference unchanged (3 unchanged) and rhiza check reports no docs/reference or docs-command failure; its only fails are .rhiza/rendered actor files for this run's untracked actor records.

## Links

### Out

- verifies [code-check-skips-reference](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-skips-reference.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
