# roadmap-umbrella-group-verified

evidence · `01M6A1WKEQ89WCCJNFQW8RD25Q` · actor [claude-code-c4f9a59d](../actor/claude-code-c4f9a59d.md)

## result

PASS

## statement

At 5ab05b16 I ran go test ./cmd/rhiza -run 'TestScripts/(roadmap|describe-golden)' -count=1: ok. The new roadmap-umbrella.txtar sets up an umbrella intent (umb), a milestone intent (m1) and change sets serving the umbrella (U), both (V), the milestone (W) and nothing (X). It asserts that roadmap --json gives group kinds milestone, backlog, umbrella in that order, that U and V each appear once (so V is in no milestone group), that the human output is '## m1 ... ## Backlog ... ## umb', and that api call roadmap gives the same bytes. describe-golden.txtar is regenerated with the new kind/groups descriptions and passes. I read the plan.go diff: planCore sends any non-draft change set serving an intent with a non-empty umbrella field (planFieldSet, also set in planCreatedNode) to that intent's umbrella group only and skips milestone and backlog placement, and those groups are appended after the backlog in intent order. make fmt-check vet test pass. make self-check fails only on stale renders for actors recorded after this commit (including this verifier session), not on the leaf. Daemon/ready treatment of umbrella groups belongs to the separate leaf umbrella-not-ready.

## Links

### Out

- verifies [roadmap-umbrella-group](../capability/rhiza/views/roadmap/roadmap-umbrella-group.md)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
