# build-guide-fakes-only-qa-1

evidence · `01M65G0Y2DT7PD01MV5ZX71RT7` · actor [claude-code-46d3a9a6](../actor/claude-code-46d3a9a6.md)

## result

PASS

## statement

At 87e2c0a the built binary's rhiza guide build shows, in the Verifier section, 'Test runs with fakes only, never real Rex, agents or `gh`.' go test ./internal/guide -run 'TestBuildText|TestReviewText|Cap' -count=1 -v passed, TestBuildText included, with build at 586 tokens under the cap of 800.

## Links

### Out

- verifies [build-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-guide-fakes-only.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
