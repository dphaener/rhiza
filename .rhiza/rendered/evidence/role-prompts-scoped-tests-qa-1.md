# role-prompts-scoped-tests-qa-1

evidence · `01M6M72RJ38B740E394RWNN7H0` · actor [claude-code-6678ddc9](../actor/claude-code-6678ddc9.md)

## result

PASS

## statement

QA at 2079b36f with the branch's binary: `rhiza workflow describe build --json` shows both prompts (builder, verifier) and `workflow describe review --json` all 8 (five checks, critic, fixer, verifier) carrying 'never the repo's full suite', with none containing go test, make check, make test or testscript; `rhiza workflow check build|review` report valid (builtin layer). go test -trimpath -count=1 ./internal/workflow/... passes, as do the prompt-printing testscripts run by name (autonomous-review, guide, arbiter-agent-tabs, run-agent-done, session-done-captures, session-run-tokens, inbox-advisor).

## Links

### Out

- verifies [role-prompts-scoped-tests](../capability/rhiza/workflows/built-in-workflows/role-prompts-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
