# release-gates

capability · `01M76B7E0P2S46PHFZ3KX010J6` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... -run TestGates with a fake clock, a fake CI status and a fixture daemon-run history file (the format tools/release defines and reads) covers each gate passing and failing, a commit missing from the history failing the soak gate, and the severe skip

## statement

tools/release check reports each gate for the candidate commit with why it passes or fails: CI green on the commit; the root daemon first ran the commit at least 24h ago and no bug request was filed since; no change set serving the finished milestone in flight and no question blocking one; a patch is at most one per 24h unless a severe: fix is in it, which skips the soak and the batching but never CI

## Links

### Out

- decided-by [m33-release-requirements](../../../decision/m33-release-requirements.md)
- decided-by [m33-batched-patches](../../../decision/m33-batched-patches.md)

### In

- verifies [release-gates-verified](../../../evidence/release-gates-verified.md)
- verifies [release-gates-qa-1](../../../evidence/release-gates-qa-1.md)
- verifies [release-gates-verified-2](../../../evidence/release-gates-verified-2.md)
- verifies [pr-498-rhiza-release-check](../../../evidence/pr-498-rhiza-release-check.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
