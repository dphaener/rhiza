# major-signoff-channel

question · `01M76B7ED4HVPNVDBZR58WS9F1`

## statement

Release tooling is for this repo only, so a major release cannot raise a Rhiza inbox item: how does the owner see and sign off a major candidate (e.g. a draft GitHub release the owner publishes, a rhiza request, or a hand-pushed tag)?

## Links

### Out

- blocks [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

### In

- answers [major-signoff-draft-release](../decision/major-signoff-draft-release.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
