# plan-acceptance-scoped-qa-1

evidence · `01M6M727AH97YN9NWQWV7M0FT0` · actor [claude-code-6678ddc9](../actor/claude-code-6678ddc9.md)

## result

PASS

## statement

QA at 2079b36f with the branch's binary: `rhiza guide plan` tells planners to write each leaf's acceptance as the tests covering that leaf, never the repo's full suite; `rhiza guide plan-review` widens missing-acceptance to an acceptance check that runs the repo's full suite; neither names go test, make check, make test or testscript; the installed rhiza-plan SKILL.md (scratch HOME) carries the rule. go test -trimpath -count=1 ./internal/guide/... passes.

## Links

### Out

- verifies [plan-acceptance-scoped](../capability/rhiza/guide/plan-acceptance-scoped.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
