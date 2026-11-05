# plan-acceptance-scoped

capability · `01M6M0F1732EDCT2C312MBA43A` · parent [guide](../guide.md)

## acceptance

A test in internal/guide asserts plan.md states the scoped-acceptance rule and plan-review.md's missing-acceptance kind covers an acceptance check that runs the repo's full suite, and that neither text contains go test, make check, make test or testscript; go test -trimpath ./internal/guide passes

## statement

The plan and plan-review guide topics, and so the rhiza-plan skill, tell planners to write each leaf's acceptance check as the tests covering that leaf, never the repo's full suite, and plan-review counts an acceptance check that runs the full suite as a missing-acceptance finding; the wording names no Go, make or testscript command

## Links

### Out

- decided-by [ci-runs-full-suite](../../../decision/ci-runs-full-suite.md)
- decided-by [scoped-tests-serve-m25-generic-rule](../../../decision/scoped-tests-serve-m25-generic-rule.md)

### In

- verifies [plan-acceptance-scoped-verified](../../../evidence/plan-acceptance-scoped-verified.md)
- verifies [plan-acceptance-scoped-qa-1](../../../evidence/plan-acceptance-scoped-qa-1.md)
- verifies [pr-259-agents-in-runs-stop-running](../../../evidence/pr-259-agents-in-runs-stop-running.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
