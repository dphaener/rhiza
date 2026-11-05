# guides-scoped-tests-verified

evidence · `01M6M6GK1SY4DBW0VNWG35PAJ4` · actor [claude-code-7acf1a93](../actor/claude-code-7acf1a93.md)

## result

PASS

## statement

At a2bd07c4, read the diff: build.md gains a Tests section telling builders and verifiers to run only the tests covering what the change touches, with the repo's own scoped commands, never the full suite, because the PR's CI runs it before merge; review.md's Checks section says the same for checks, critics and fixers. TestBuildText and TestReviewText assert those phrases and that neither body contains go test, make check, make test or testscript; grep of both texts finds none of the four. go test -trimpath -count=1 ./internal/guide passes (TestBuildText and TestReviewText PASS under -v); go vet and gofmt -l on internal/guide are clean.

## Links

### Out

- verifies [guides-scoped-tests](../capability/rhiza/guide/guides-scoped-tests.md)

## History

- `01M6M01GYQHR9VMCZ54W5THARK` Agents in runs stop running the full test suite locally: builders, verifiers, review checks, critics and fixers run only the tests of the packages and scripts their change touches (scoped go test with -trimpath, the affected testscripts, make self-check), and no step runs make check/make test or a bare go test ./... — the PR's GitHub CI runs the full suite before every merge and the ship workflow already waits for it. Owner, 2026-11-04: "Yes the test suites make sense, since we're going to run full CI anyway on GitHub before merge." (context: the VM's host CPU is over-committed, ~35% steal; every agent ran a full suite). Update the built-in workflows' role prompts, the rhiza-build/rhiza-review skills and guide topics, and CLAUDE.md's guidance, and keep the leaf acceptance checks scoped
