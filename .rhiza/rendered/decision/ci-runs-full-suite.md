# ci-runs-full-suite

decision · `01M6M0ECQ2TXFF5Y2V9J59HNX8`

## quote

Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge.

## statement

Agents in runs run only the tests covering what their change touches and never the repo's full suite; the PR's GitHub CI runs the full suite before every merge

## Links

### Out

- derived-from [scoped-tests-plan](../source/scoped-tests-plan.md)

### In

- relates-to [no-ci-as-debugger](no-ci-as-debugger.md) — Refines ci-runs-full-suite for unstable CI; does not replace it
- decided-by [guides-scoped-tests](../capability/rhiza/guide/guides-scoped-tests.md)
- decided-by [plan-acceptance-scoped](../capability/rhiza/guide/plan-acceptance-scoped.md)
- decided-by [role-prompts-scoped-tests](../capability/rhiza/workflows/built-in-workflows/role-prompts-scoped-tests.md)
- decided-by [claude-md-scoped-tests](../capability/rhiza/integrity/self-check/claude-md-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
