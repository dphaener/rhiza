# release-next-version-verified-3

evidence · `01M787QA4Y3QMD37VW2VY7PANM` · actor [claude-code-270bbc78](../actor/claude-code-270bbc78.md)

## result

PASS

## statement

At 89df6c69e (fixer commit: unlandable servers keep no milestone open) ran go test -trimpath -count=1 ./tools/release/... -run 'TestNextVersion|TestMilestone' -v: TestNextVersion, TestNextVersionFromLogs and the new TestMilestoneIgnoresUnlandableServers pass; go vet clean. Read the diff: ReadHistory now skips serves from change sets whose withdraw/reject/supersede entries end them (entry type names match internal/entry TypeWithdraw/TypeReject/TypeSupersede; reopen only follows a landing per the lifecycle table), serves links retired by link.retire (unless the retire op is dropped) and link.create ops a body.drop retracts; op field names id/link/drop match real log entries. go run ./tools/release check on this repo's log reports v0.2.0 minor since v0.1.0 (pre-v1 breaking is minor) with milestone-32 finished.

## Links

### Out

- verifies [release-next-version](../capability/rhiza/distribution/release-next-version.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
