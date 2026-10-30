# review-guide-fakes-only-qa-1

evidence · `01M65G130JTAM3K2K0J6EDKS52` · actor [claude-code-46d3a9a6](../actor/claude-code-46d3a9a6.md)

## result

PASS

## statement

At 87e2c0a the built binary's rhiza guide review shows 'Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`.' at the head of the Checks list, and the qa line still says hands-on with the built binary, fakes alone do not count. go test ./internal/guide -run 'TestBuildText|TestReviewText|Cap' -count=1 -v passed, TestReviewText included, with review at 728 tokens under the cap of 800.

## Links

### Out

- verifies [review-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-guide-fakes-only.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
