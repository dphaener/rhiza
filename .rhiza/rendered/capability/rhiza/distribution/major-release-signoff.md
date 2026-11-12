# major-release-signoff

capability · `01M76B7EBC8W8NQW1GXBYYB3Z0` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... -run TestMajor holds a major until the sign-off is given and then cuts it once

## statement

A major candidate is never tagged by the tool: it waits for the owner sign-off, then tags and pushes as for a minor

## Links

### Out

- decided-by [m33-auto-minor-patch](../../../decision/m33-auto-minor-patch.md)
- decided-by [release-annotated-tags-only](../../../decision/release-annotated-tags-only.md)

### In

- blocks [major-signoff-channel](../../../question/major-signoff-channel.md)
- verifies [major-release-signoff-verified](../../../evidence/major-release-signoff-verified.md)
- verifies [major-release-signoff-qa-1](../../../evidence/major-release-signoff-qa-1.md)
- verifies [major-release-signoff-verified-2](../../../evidence/major-release-signoff-verified-2.md)
- verifies [major-release-signoff-qa-2](../../../evidence/major-release-signoff-qa-2.md)
- verifies [major-release-signoff-verified-3](../../../evidence/major-release-signoff-verified-3.md)
- verifies [pr-506-automatic-releases](../../../evidence/pr-506-automatic-releases.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
