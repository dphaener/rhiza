# web-run-detail-verified-2

evidence · `01M6Q9D45Y1WE4JM0BEWTTRREQ` · actor [claude-code-ebae0403](../actor/claude-code-ebae0403.md)

## result

PASS

## statement

At ea99839d, go test -trimpath -count=1 -run 'TestRunPage' ./internal/web/ passes: TestRunPagePlanRun renders /runs/<id> for a fixture plan run and asserts steps in order (interview, draft, gate) with done/failed/running outcomes, durations and times, and the sessions with status ('finished · 1115 tokens', now counting cache_creation and cache_read, and 'running'); TestRunPageReviewRun asserts the review summary's check, findings and severities in order (blocker Nil dereference a.go:3, nit Spelling); TestRunPageUnknownIs404 gets 404 for an unknown run; TestRunPageThroughReadOnlyAPI serves a real run journal through api.ReadOnlyHandler and gets 200 with the title on /runs/run1 and a rendered fragment. TestReadOnlyMarker and TestRunStatusAndResumeRefusals in internal/api pass, and runRunStatus only reads the journal, session store and ship state, so its ReadOnly mark is sound. The diff's token struct tags match api.SessionTokens' json names. The full internal/api package timed out under machine load (load avg 12-19) on the unrelated TestChangeShowIndexedMatchesFoldGenerated; that is left to CI.

## Links

### Out

- verifies [web-run-detail](../capability/rhiza/web-monitor/web-run-detail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
