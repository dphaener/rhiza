# role-prompts-scoped-tests-verified

evidence · `01M6M6Q2FE636WNWGN3GX1PAYP` · actor [claude-code-4dda77b1](../actor/claude-code-4dda77b1.md)

## result

PASS

## statement

Verifier claude-code-4dda77b1 at 4d79c2be read the diff: build.toml (builder, verifier) and review.toml (5 checks, critic, fixer, verifier) prompts each append 'Run only the tests covering what your change touches, never the repo's full suite: the PR's CI runs it before merge.' and name no go test, make check, make test or testscript. New TestBuiltinRolePromptsScopeTests in internal/workflow parses both built-ins via Builtin/Parse, collects every builder/verifier/critic/fixer/check prompt, and asserts the rule is present and none of the four commands appear. go test -trimpath -count=1 ./internal/workflow: ok (test PASS under -v). No cmd/rhiza testscript golden prints these prompts (grep for their text finds none); as a smoke check go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(build-evidence-commit|run-ship-signoff|review-signoff-changes|inbox-signoff|autonomous-review|run-ship-change|guide-install|describe-golden)$': ok. The commit's only .rhiza/log files are the builder's own actor record.

## Links

### Out

- verifies [role-prompts-scoped-tests](../capability/rhiza/workflows/built-in-workflows/role-prompts-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
