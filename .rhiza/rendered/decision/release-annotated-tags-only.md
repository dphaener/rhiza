# release-annotated-tags-only

decision · `01M78RMAZBWX0SZA38Y1976YC5`

## options

```json
[
  "describe annotated tags only (no --tags)",
  "describe --tags and skip non-annotated v* tags"
]
```

## statement

Only an annotated v* tag is a release: evaluate (tools/release/next.go) and breakingAgainstLastTag (tools/release/breaking.go) run git describe without --tags, so the lightweight vX.0.0 tag that publishing a major's draft GitHub release makes is the sign-off, not the last release, and Cut still tags the signed-off major with its annotated tag. Recorded by the review cycle 3 fixer from what commit 4d496116f does; no owner words.

## Links

### Out

- answers [major-signoff-lightweight-tag](../question/major-signoff-lightweight-tag.md)

### In

- decided-by [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
