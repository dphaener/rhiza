# major-release-signoff-qa-1

evidence · `01M78QJP65N4GJ15T6YGRGZQ9F` · actor [claude-code-ddf76b91](../actor/claude-code-ddf76b91.md)

## result

FAIL

## statement

QA at 903604b3b, built tools/release binary, real git with scratch clone and bare origin, fake gh. With v1.0.0 tagged and a breaking difference since, cut drafted one GitHub release v2.0.0 at the candidate (gh release create --draft --target <candidate>) and pushed no tag; with the draft present a second run created no new draft and pushed nothing. FAIL: after simulating the owner's sign-off (lightweight v2.0.0 at the candidate on origin, gh release list reporting it published) cut's fetch --tags brings the lightweight tag in, git describe --tags then takes v2.0.0 as the last release and cut prints 'none: no release since v2.0.0' and exits 0: the tag on origin stays a lightweight commit tag, so release.yml skips GoReleaser and the signed-off major is never cut.

## Links

### Out

- verifies [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
