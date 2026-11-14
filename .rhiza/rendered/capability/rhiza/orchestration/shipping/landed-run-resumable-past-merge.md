# landed-run-resumable-past-merge

capability · `01M7AP37F8KV6E7D4SWA7XTBW2` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./internal/run -run TestPastLanding shows the method true for a run that visited land, merge, ship-change/land or ship-change/merge, and false for one that visited neither (including a step such as ship-change/merge-check or premerge whose last segment is not land or merge); go test -trimpath ./internal/daemon -run TestResumableOnlyAcceptedOrLandedPastLand shows an interrupted run on a landed change set that visited merge but not land picked, one that visited ship-change/merge picked, and one that visited neither skipped; go test -trimpath ./internal/api -run TestRunResumeLandedPastMerge shows rhiza run resume accepting such a run (bare merge and ship-change/merge) and still refusing a landed run that visited neither; go test -trimpath ./cmd/rhiza -run TestScripts/inbox-landed-resumable passes with fixture runs that visited merge, and ship-change/merge, but not land listed and not skipped. Touches internal/run/journal.go, internal/run/journal_test.go, internal/api/run.go, internal/api/run_test.go, internal/daemon/pick.go, internal/daemon/pick_test.go, internal/api/inbox.go and cmd/rhiza/testdata/script/inbox-landed-resumable.txtar.

## statement

A non-plan run on a landed change set counts as past its landing when it entered a land step or a merge step, not only a land step, so a conductor-default run (which lands on the branch and has no land step) interrupted after its merge (main's CI, close-out, ship, deploy-check) is resumable: rhiza run resume lets it through instead of refusing it as no longer accepted, the daemon's resumable pick takes it, and the inbox lists it instead of leaving it out as gone. A step counts as land or merge when its name is land or merge, or its last segment after the final "/" is, so an included step (ship-change/merge in this repo's ship-rhiza and bugfix runs) counts too. One run.State method decides it for all three; a landed run that entered neither stays refused, skipped and gone as today.

## Links

### Out

- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)

### In

- verifies [landed-run-resumable-past-merge-verified](../../../../evidence/landed-run-resumable-past-merge-verified.md)
- verifies [landed-run-resumable-past-merge-qa-2](../../../../evidence/landed-run-resumable-past-merge-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
