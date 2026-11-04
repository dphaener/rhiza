# question-c2063c96-1

question · `01M6M06NCT2Z31K6CF81ASAW2T`

## session

c2063c96-e994-4af5-babc-01a1a800e15c

## statement

Change set 'Agents in runs stop running' (scoped tests in runs; CI runs the full suite) serves no intent yet. Which landed intent should it serve? / The build/review/plan guide topics (and so the rhiza-build/-review/-plan skills) are product text used in every repo, Nebari included. How specific should their scoped-test rule be?

## Links

### In

- answers [scoped-tests-serve-m25-generic-rule](../decision/scoped-tests-serve-m25-generic-rule.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
