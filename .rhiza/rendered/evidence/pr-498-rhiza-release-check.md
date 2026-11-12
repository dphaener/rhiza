# pr-498-rhiza-release-check

evidence · `01M788FDEJ417R02Y9B8H47F5A` · actor [dphaener](../actor/dphaener.md)

## head

39138ed83e6dfc50c025418fcb0856a304e8ae93

## pr

https://github.com/dphaener/rhiza/pull/498

## result

13 checks passed: changes SUCCESS, snapshot SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #498 (https://github.com/dphaener/rhiza/pull/498) passed CI at 39138ed

## status

verified

## Links

### Out

- verifies [release-next-version](../capability/rhiza/distribution/release-next-version.md)
- verifies [release-gates](../capability/rhiza/distribution/release-gates.md)
- verifies [release-notes](../capability/rhiza/distribution/release-notes.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
