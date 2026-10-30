# review-guide-fakes-only-verified

evidence · `01M65FXM0JG01MA0H2ZSD8B64H` · actor [claude-code-717f8f21](../actor/claude-code-717f8f21.md)

## result

PASS

## statement

At 42217dd, read the diff: internal/guide/text/review.md adds 'Every check except `qa` tests runs with fakes only, never real Rex, agents or `gh`.' at the head of the Checks list (covering compliance, quality, scope and model), and the existing qa line still says hands-on with the built binary, fakes alone do not count. Ran go test ./internal/guide -run 'TestReviewText|Cap' -v: TestReviewText PASS, its new fakes-only group pins that sentence (naming Rex, agents and gh) and its five-checks group pins 'built binary' and 'fakes alone do not count'; TestCaps PASS with review at 728 of 800 tokens (TopicCap). go test ./... passes with no failures and gofmt -l is empty. make check's self-check only reports stale .rhiza/rendered files from the run's actor records, not from this leaf.

## Links

### Out

- verifies [review-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-guide-fakes-only.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
