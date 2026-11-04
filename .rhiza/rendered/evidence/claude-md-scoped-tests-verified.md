# claude-md-scoped-tests-verified

evidence · `01M6M69TN3F8P9HPV11V1827K4` · actor [claude-code-c5bebbed](../actor/claude-code-c5bebbed.md)

## result

PASS

## statement

Verifier claude-code-c5bebbed at f99ccf13: read the diff (CLAUDE.md replaces the unqualified 'Run `make check` before calling any change done.' with the agents-in-runs rule: never run make check, make test or a bare go test ./..., since GitHub CI runs make check before every merge; scoped commands go test -trimpath ./internal/<pkg>/..., go test -trimpath ./cmd/rhiza -run TestScripts/<name>, make self-check; make check kept for people outside a run). Ran `go test -trimpath ./cmd/rhiza -run 'TestClaudeMDScopedTests$' -count=1 -v`: PASS. Mutation: same test in an archive of HEAD with CLAUDE.md swapped for HEAD~1's version FAILS on all four missing phrases and on the still-present unqualified line. gofmt clean.

## Links

### Out

- verifies [claude-md-scoped-tests](../capability/rhiza/integrity/self-check/claude-md-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
