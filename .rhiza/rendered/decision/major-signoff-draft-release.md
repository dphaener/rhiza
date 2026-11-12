# major-signoff-draft-release

decision · `01M76NPP2GYFJRH1PWXEB3FWNZ`

## question

Release tooling is for this repo only, so a major release cannot raise a Rhiza inbox item: how does the owner see and sign off a major candidate (e.g. a draft GitHub release the owner publishes, a rhiza request, or a hand-pushed tag)?

## statement

A major release candidate is prepared as a draft GitHub release with notes; the owner publishing it is the sign-off, after which the tool tags and pushes

## wording

Draft GitHub release: the tool prepares a draft GitHub release with notes, publishing it is the owner sign-off, and the tool then tags and pushes; picked by the owner

## Links

### Out

- answers [major-signoff-channel](../question/major-signoff-channel.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
