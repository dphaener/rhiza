# claude-md-scoped-tests

capability · `01M6M0FK7XV5HMP0A70DZFF188` · parent [self-check](../self-check.md)

## acceptance

A test reads CLAUDE.md and finds go test -trimpath, a TestScripts/<name> run and make self-check as the agents' scoped commands, and the rule that agents in runs never run make check, make test or a bare go test ./...; it no longer finds the unqualified line 'Run `make check` before calling any change done.'; that test passes with go test -trimpath on its package

## statement

This repo's CLAUDE.md tells agents in runs to test only what their change touches, with go test -trimpath on the touched packages, the affected cmd/rhiza testscripts run by name, and make self-check, and never make check, make test or a bare go test ./..., since the PR's GitHub CI runs make check before every merge

## Links

### Out

- decided-by [ci-runs-full-suite](../../../../decision/ci-runs-full-suite.md)
- decided-by [scoped-tests-serve-m25-generic-rule](../../../../decision/scoped-tests-serve-m25-generic-rule.md)

### In

- verifies [claude-md-scoped-tests-verified](../../../../evidence/claude-md-scoped-tests-verified.md)
- verifies [claude-md-scoped-tests-qa-1](../../../../evidence/claude-md-scoped-tests-qa-1.md)
- verifies [pr-259-agents-in-runs-stop-running](../../../../evidence/pr-259-agents-in-runs-stop-running.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
