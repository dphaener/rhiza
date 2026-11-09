# code-check-leaves-reference-to-self-check

decision · `01M6ZTCHFZNN7V7KTVN25KBDPW`

## reasoning

m30-reference says the docs/reference files are generated from the code and make check fails when they are stale; #412 (pr-412-generated-reference evidence) served it, and this change set's intent names that regression, so the fix belongs to the same intent rather than an older run-reliability intent like m4-runs-hold-up. Both rules compare the repo against the binary's own registry (docsFiles() and Operations() in internal/api/check.go), so they are only correct when the binary is the branch's build; the daemon's in-process check is the installed binary, which is why a run that edits CLI help, settings or the workflow format, or adds a documented command, can never pass review/code until merged and installed. Building the worktree's binary inside the code check was weighed and rejected: the code check is generic across projects under the root, only Rhiza's own repo has a rhiza-generated docs/reference, and a go build per check adds Go-toolchain and CPU cost on the shared VM. make self-check already builds this commit and runs in CI (make check) before every merge, so turning the two rules off in the code check loses no coverage on the merge path, and an explicit request option keeps the default rhiza check unchanged for people and CI.

## statement

Change set 01M6ZT7 serves intent m30-reference (01M6W5BP8JF2): it repairs a regression from #412, the change set that served that intent. The run code check (internal/api/run_codecheck.go checkIn/runCheck) stops judging the tool's own reference: it runs rhiza check with the generated-reference comparison (DocsWant/Docs, the render failure 'differs from what rhiza docs writes') and the docs-command rule (Commands) turned off, through an explicit CheckRequest option the code check sets, so plain rhiza check and make self-check keep both. Those two rules are judged only by make self-check, which builds the branch and runs in the PR's CI before every merge. The code check does not build the run worktree's binary.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-e5dcdc72-1](../question/question-e5dcdc72-1.md)
- derived-from [plan-code-check-reference](../source/plan-code-check-reference.md)

### In

- decided-by [code-check-skips-reference](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-skips-reference.md)

## History

- `01M6ZT7VZZWD3R8F5HVWXRN1TP` Fix Every run whose change edits CLI help, settings or the workflow format now fails at review/code. The code check runs rhiza check in-process (internal/api/run_codecheck.go:132 runCheck) with the installed binary, and since #412 (generated reference) check compares docs/reference/ with what that binary's rhiza docs writes, not the branch's build. Run 01M6ZRM3QDK5VYBER6DYZDDR65 (change set 01M6ZQ588DTH): 'rhiza check: fail render docs/reference/cli.md: differs from what rhiza docs writes'. In its worktree the branch build's check passes ('3 unchanged' from its rhiza docs) while ~/bin/rhiza check fails the same way. Retry can never pass until the change is merged and installed. Fix: the code check judges docs/reference with the run worktree's own build (as make self-check does), or leaves that one check to make self-check/CI.
