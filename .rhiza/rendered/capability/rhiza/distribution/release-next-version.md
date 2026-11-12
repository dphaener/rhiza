# release-next-version

capability · `01M76B7E06G9AJHZSV3GSEA005` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... -run TestNextVersion over fixture logs covers each kind, the pre-v1 breaking-is-minor rule and no change since the tag

## statement

tools/release check computes the next version and its kind since the last v* tag: major from v1.0.0 on when a change is marked or detected breaking (before v1 a breaking change makes a minor), minor when a milestone (an intent with a milestone field) finished, that is every change set serving it or an intent under it landed; otherwise patch when a bug fix landed, none otherwise. A patch is cut from main, so it ships whatever landed beside the fix, an unfinished milestone's features included

## Links

### Out

- decided-by [m33-milestone-is-minor](../../../decision/m33-milestone-is-minor.md)
- decided-by [m33-breaking-is-major](../../../decision/m33-breaking-is-major.md)

### In

- verifies [release-next-version-verified](../../../evidence/release-next-version-verified.md)
- verifies [release-next-version-qa-1](../../../evidence/release-next-version-qa-1.md)
- verifies [release-next-version-verified-2](../../../evidence/release-next-version-verified-2.md)
- verifies [release-next-version-verified-3](../../../evidence/release-next-version-verified-3.md)
- verifies [pr-498-rhiza-release-check](../../../evidence/pr-498-rhiza-release-check.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
