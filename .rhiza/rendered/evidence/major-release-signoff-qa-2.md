# major-release-signoff-qa-2

evidence · `01M78RFJ99YDE9AZTQGZ621V7Q` · actor [claude-code-03513d81](../actor/claude-code-03513d81.md)

## result

PASS

## statement

QA at 8f0616278 with the built tools/release binary, real git with a scratch clone and bare origin, and a fake gh. The last release was an annotated v1.0.0, with breaking changes since. The first cut drafted one GitHub release, v2.0.0 at the candidate (gh release create --draft --target), and pushed no tag. With the draft reported, a second cut created no new draft and pushed nothing. Next I simulated the owner's sign-off: a lightweight v2.0.0 at the candidate on origin, with gh release list reporting it published. cut printed 'v2.0.0 is signed off by the owner' and 'tagged and pushed v2.0.0'. Origin's v2.0.0 is now an annotated tag object at the candidate with all 5 '## ' headings. That fixes the cycle-1 failure, where the lightweight tag was taken as the last release. A repeat cut printed 'none: no release since v2.0.0' and left origin's tag unchanged. Exactly one draft was created in total.

## Links

### Out

- verifies [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
