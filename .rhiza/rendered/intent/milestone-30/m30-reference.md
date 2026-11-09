# m30-reference

intent · `01M6W5BP8JF2CV5FHPZAZNR0RP` · parent [milestone-30](../milestone-30.md)

## statement

The CLI, rhiza.toml settings and workflow file references in docs/reference/ are generated from the code, and make check fails when they are stale.

## Links

### Out

- decided-by [m30-doc-set](../../decision/m30-doc-set.md)
- decided-by [m30-generated-reference](../../decision/m30-generated-reference.md)

### In

- serves change set `01M6W5E3JQXD62WN4RHG546KRX` (Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist)
- serves change set `01M6ZT7VZZWD3R8F5HVWXRN1TP` (Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.)

## History

- `01M6W5A2R4BXDGQGHNWGZNB3XY` Milestone 30: Rhiza is ready to share: install, README and docs, generated reference, public under MIT
