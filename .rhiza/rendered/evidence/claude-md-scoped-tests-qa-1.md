# claude-md-scoped-tests-qa-1

evidence · `01M6M710NNSRFDF46EHB3VA70M` · actor [claude-code-6678ddc9](../actor/claude-code-6678ddc9.md)

## result

PASS

## statement

QA at 2079b36f: read CLAUDE.md in the built tree; it names go test -trimpath ./internal/<pkg>/..., go test -trimpath ./cmd/rhiza -run TestScripts/<name> and make self-check as the agents' scoped commands, says agents in runs never run make check, make test or a bare go test ./... because GitHub CI runs make check before merge, and keeps make check only for a person outside a run. go test -trimpath -count=1 ./cmd/rhiza -run 'TestClaudeMDScopedTests$' passes.

## Links

### Out

- verifies [claude-md-scoped-tests](../capability/rhiza/integrity/self-check/claude-md-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
