# web-roadmap-page-verified-2

evidence · `01M6Q8FC3KQZS1JQQ3T0ND37AK` · actor [claude-code-65b4429b](../actor/claude-code-65b4429b.md)

## result

PASS

## statement

At 81427c0d, go test -trimpath -count=1 -run TestRoadmap -v ./internal/web/ passes (TestRoadmapPageMatchesTheOperation, TestRoadmapPageEmptyAndNoAPI), as does the whole internal/web package; go vet and gofmt are clean. The test builds a fixture with two milestones, a backlog change set and an umbrella, checks the roadmap operation returns milestone, milestone, backlog, umbrella, then asserts /roadmap shows the titles and change sets in that order with the umbrella labelled Umbrella after the backlog, a /changes/<id> link and a state pill for each change set, no IDs or handles, a kicker counting 2 milestones (umbrella excluded), and a live fragment that picks up a newly planned change set. Read the diff: RoadmapView now keeps the operation's groups in order instead of splitting milestones and backlog, which fixes the QA finding that the umbrella group showed as a third milestone before the backlog; the leaf statement is amended to drop the short ID and say rows show the intent. make self-check (RHIZA_ROOT unset) fails only on stale .rhiza/rendered files, which are rendered at land.

## Links

### Out

- verifies [web-roadmap-page](../capability/rhiza/web-monitor/web-roadmap-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
