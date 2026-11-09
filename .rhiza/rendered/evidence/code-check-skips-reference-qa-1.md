# code-check-skips-reference-qa-1

evidence · `01M6ZWGRS0X2053Z2WZ9FCEVJN` · actor [claude-code-7a2a518d](../actor/claude-code-7a2a518d.md)

## result

PASS

## statement

QA with the branch build (go build -trimpath ./cmd/rhiza) in a scratch root and repo, hermetic HOME: started a run on a wait-step workflow, committed leaf alpha in the run worktree with verifier Evidence; rhiza run code-check passed. Then ran rhiza docs, appended a stray line to docs/reference/cli.md and wrote a README naming `rhiza no-such-command`: plain rhiza check failed with 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and 'fail docs-command README.md', exit 1; rhiza check --skip-reference was ok, exit 0, and printed the same bytes as rhiza api call check '{"skip_reference":true}'; rhiza run code-check still passed with neither among its causes. A stale .rhiza/rendered/index.md still failed code-check with the render cause. go test ./internal/api -run TestCheck and TestScripts/review-code-check and describe-golden pass.

## Links

### Out

- verifies [code-check-skips-reference](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-skips-reference.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
