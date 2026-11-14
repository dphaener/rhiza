# ship-one-ci-round-verified

evidence · `01M7C5P7RXBZ5DYKWCRAEM8143` · actor [claude-code-e737a3c4](../actor/claude-code-e737a3c4.md)

## result

PASS

## statement

At 7f1701a77 re-ran the Checks block: go test -trimpath ./cmd/rhiza -run TestScripts/ship-change-one-ci (exit 0, the acceptance command), go test -trimpath ./internal/api -run 'TestReviewJournaled|TestShipOpenPR' (exit 0), go test -trimpath ./internal/workflow (exit 0), rhiza ship generate and rhiza render with a scratch build of this commit (exit 0, tree clean); go vet and gofmt clean. Read the diff: the script uses a fake gh, a fake rex and a fake conductor that lands and renders on the branch. On the happy path it asserts one wait-ci visit, one push of cs/ship-it in the origin's pre-receive log, one pr merge, no land step and the change set landed with main's rendered index. When main moves at merge it asserts merge visit 1 fails, conduct visit 2 runs, wait-ci visit 2 passes, two pushes, one merge, one hand-back and the change set landed. Beyond the txtar: ship_openpr.go reviewJournaled now counts a conductor check's rounds (unit-tested by TestReviewJournaledConductorRounds); this is allowed under internal/api/ship*.go. The builtin ship-change.toml conduct step gains a handback template, which the commit explains and which is required: stepRunner.Handback errors when a step has no handback template, and the main-moved path exercises it. No stubs or TODOs.

## Links

### Out

- verifies [ship-one-ci-round](../capability/rhiza/orchestration/shipping/ship-one-ci-round.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
