# code-check-skips-reference-verified

evidence · `01M6ZWAV649GHJ3TKP9VCK4PYF` · actor [claude-code-76a116e9](../actor/claude-code-76a116e9.md)

## result

PASS

## statement

At 94281d07 I read the diff: CheckRequest gains SkipReference (json skip_reference), runCheck leaves DocsWant, Commands and Docs unset when it is true, run code-check's checkIn passes SkipReference: true, and the only other caller (health) and plain rhiza check keep the default false, so both rules stay for rhiza check and make self-check. go test -trimpath ./internal/api/... passed, including TestCheckSkipReference (default check has 2 fail findings, on docs/reference/cli.md and README.md; skip_reference has 0). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(review-code-check|describe-golden)$' passed: review-code-check makes docs/reference stale and a README naming rhiza no-such-command, sees rhiza check fail naming both, then rhiza run code-check pass with exit 0 and neither in its output; describe-golden includes skip_reference.

## Links

### Out

- verifies [code-check-skips-reference](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-skips-reference.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
