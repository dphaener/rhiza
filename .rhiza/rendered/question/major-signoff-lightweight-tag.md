# major-signoff-lightweight-tag

question · `01M78QVJCD10S11S58V9NR3E8K`

## statement

Found with real git while fixing review cycle 2 (not fixed: outside the findings). Once the owner publishes a major's draft release, GitHub makes a lightweight tag vX.0.0 at the candidate. Cut's `git fetch --tags` brings it here, evaluate's `git describe --tags` (next.go:355) then takes it as the last release at HEAD, and Cut prints 'none: no release since vX.0.0' instead of replacing it with the annotated tag. TestMajor's fake git does not model this. Should evaluate ignore a lightweight tag (describe without --tags, or skip non-annotated v* tags) so the signed-off major still gets its annotated tag?

## Links

### Out

- relates-to [auto-release](../capability/rhiza/distribution/auto-release.md) — a signed-off major may never get its annotated tag

### In

- answers [release-annotated-tags-only](../decision/release-annotated-tags-only.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
