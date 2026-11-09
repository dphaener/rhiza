# pr-417-fix-every-run-whose-change

evidence · `01M700QHCW7KN3W04ZA4TFR695` · actor [dphaener](../actor/dphaener.md)

## head

719df8f2b06d87b94e9ff78063f7e9d246156830

## pr

https://github.com/dphaener/rhiza/pull/417

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #417 (https://github.com/dphaener/rhiza/pull/417) passed CI at 719df8f

## status

verified

## Links

### Out

- verifies [code-check-skips-reference](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-skips-reference.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
