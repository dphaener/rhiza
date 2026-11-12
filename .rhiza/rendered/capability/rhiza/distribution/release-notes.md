# release-notes

capability · `01M76B7E153HS76CMR3ZHVPTM6` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... -run TestNotes compares notes from a fixture log with a golden file

## statement

tools/release check drafts the release notes from the change sets landed since the last tag, grouped by milestone, with bug fixes and breaking changes (and their reasons) called out

## Links

### Out

- decided-by [m33-release-requirements](../../../decision/m33-release-requirements.md)

### In

- verifies [release-notes-verified](../../../evidence/release-notes-verified.md)
- verifies [release-notes-qa-1](../../../evidence/release-notes-qa-1.md)
- verifies [pr-498-rhiza-release-check](../../../evidence/pr-498-rhiza-release-check.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
