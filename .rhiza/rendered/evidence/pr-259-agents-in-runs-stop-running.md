# pr-259-agents-in-runs-stop-running

evidence · `01M6M7VE6EDCT05HYGH9X1YJ26` · actor [dphaener](../actor/dphaener.md)

## head

23a21a182ecd6a3a0a1064c448f150ce0d5fd1e4

## pr

https://github.com/dphaener/rhiza/pull/259

## result

1 check passed: check SUCCESS

## statement

Pull request #259 (https://github.com/dphaener/rhiza/pull/259) passed CI at 23a21a1

## status

verified

## Links

### Out

- verifies [guides-scoped-tests](../capability/rhiza/guide/guides-scoped-tests.md)
- verifies [plan-acceptance-scoped](../capability/rhiza/guide/plan-acceptance-scoped.md)
- verifies [role-prompts-scoped-tests](../capability/rhiza/workflows/built-in-workflows/role-prompts-scoped-tests.md)
- verifies [claude-md-scoped-tests](../capability/rhiza/integrity/self-check/claude-md-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
