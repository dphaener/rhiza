# plan-acceptance-scoped-verified

evidence · `01M6M6JWPQFXM40AB6CKHR6705` · actor [claude-code-85a52809](../actor/claude-code-85a52809.md)

## result

PASS

## statement

At 2f56cfa3, ran go test -trimpath -count=1 ./internal/guide/... (ok) and -run TestPlanAcceptanceScoped -v (PASS). The new test asserts plan.md says to write each leaf's acceptance as 'the tests covering that leaf, never the repo's full suite', asserts plan-review.md's missing-acceptance kind covers one that 'runs the repo's full suite' rather than the 'tests covering that leaf', and fails if either body contains go test, make check, make test or testscript; grep of both texts finds none of those. plan.md is the rhiza-plan skill body (Lookup(plan).Skill), so the skill carries the rule. Diff read: only the two guide texts, the test and the builder's actor record.

## Links

### Out

- verifies [plan-acceptance-scoped](../capability/rhiza/guide/plan-acceptance-scoped.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
