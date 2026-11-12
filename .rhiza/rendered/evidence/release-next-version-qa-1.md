# release-next-version-qa-1

evidence · `01M786MRK1KX1XSCVB23D3W2B1` · actor [claude-code-c6f5cfbe](../actor/claude-code-c6f5cfbe.md)

## result

PASS

## statement

QA at 8ff82f55d: ran the built `release check` by hand. On this repo it reports v0.2.0 minor since v0.1.0 (two detected setting breaks, so breaking before v1 is a minor). On a scratch clone tagged v0.2.0: only a landed bug fix gives v0.2.1 patch; after tagging v0.2.1 it gives v0.2.2; a change set with a 'breaking: removes --old' note gives v0.3.0 minor (pre-v1); tagging v1.0.0 makes it v2.0.0 major; with nothing landed since the tag it reports 'none: no release since v0.2.0'.

## Links

### Out

- verifies [release-next-version](../capability/rhiza/distribution/release-next-version.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
