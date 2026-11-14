# landed-run-resumable-past-merge-verified

evidence · `01M7C39G1QZFEB64HYC9BFPERB` · actor [claude-code-220ebd7f](../actor/claude-code-220ebd7f.md)

## result

PASS

## statement

At f2ea977fa re-ran the Checks block and acceptance: go test -trimpath ./internal/run -run TestPastLanding, ./internal/daemon -run TestResumableOnlyAcceptedOrLandedPastLand, ./internal/api -run 'TestRunResume|TestInbox' and -run TestRunResumeLandedPastMerge (PASS), ./cmd/rhiza -run TestScripts/inbox-landed-resumable all exit 0; rhiza ship generate and rhiza render exit 0. Diff adds State.PastLanding (land/merge by name or last segment after /), used by run resume, daemon Resumable and inbox gone; tests cover merge, ship-change/merge, ship-change/land picked/resumed/not skipped and ship-change/merge-check, premerge, landing refused. The api test sits in dispatch_request_test.go beside the other TestRunResume tests rather than run_test.go. No stubs or TODOs, nothing outside the leaf.

## Links

### Out

- verifies [landed-run-resumable-past-merge](../capability/rhiza/orchestration/shipping/landed-run-resumable-past-merge.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
