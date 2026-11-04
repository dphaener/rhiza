# scoped-tests-serve-m25-generic-rule

decision · `01M6M0BQGFN47ZDN4APT1K9BZ0`

## reasoning

Intent: milestone-25 is the daemon emptying the backlog on its own, and the landed change sets about the cost of tests and runs on the VM (01M6JJMK3M50 test suite cheap in many worktrees, 01M6HR3VBT4D timing budgets in quiet, 01M6KV89S5MK cheap inbox, 01M6HW7MG0ZV github runner) all serve it; this one removes the same CPU contention (~35% steal) that slows the daemon's runs, so it fits there, not m24 (VM setup, done) or m26 (Nebari import). Specificity: internal/guide/text/build.md and review.md are generic today and mention no make or go commands, and the question itself notes they ship to every repo including Nebari, which is not a Go/make repo; a Go-specific rule there would be wrong advice elsewhere. The owner's reason ('we're going to run full CI anyway on GitHub before merge') is repo-independent, so the generic rule carries it, while the repo-specific commands belong where repo specifics already live: CLAUDE.md and the acceptance checks, which the change set's intent already lists.

## statement

The change set serves milestone-25 (01M6HKVRE1B7RPP12EN6J2GGXC). The build/review/plan guide topics and the built-in workflows' role prompts state the scoped-test rule repo-agnostically: run only the tests covering what your change touches, using the repo's own scoped commands, and never the repo's full suite, because the PR's CI runs it before merge. They name no Go, make or testscript commands; this repo's concrete commands (go test -trimpath on the touched packages, the affected testscripts, make self-check, never make check/make test/bare go test ./...) go in its CLAUDE.md and in the leaves' acceptance checks.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c2063c96-1](../question/question-c2063c96-1.md)
- derived-from [scoped-tests-plan](../source/scoped-tests-plan.md)

### In

- decided-by [guides-scoped-tests](../capability/rhiza/guide/guides-scoped-tests.md)
- decided-by [plan-acceptance-scoped](../capability/rhiza/guide/plan-acceptance-scoped.md)
- decided-by [role-prompts-scoped-tests](../capability/rhiza/workflows/built-in-workflows/role-prompts-scoped-tests.md)
- decided-by [claude-md-scoped-tests](../capability/rhiza/integrity/self-check/claude-md-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
