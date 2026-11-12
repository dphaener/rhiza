# pr-506-automatic-releases

evidence · `01M78W1MG81K9H0SN7QB07D3Q8` · actor [dphaener](../actor/dphaener.md)

## head

1dad0c40c5c5f6aadfd947b5b604f34af87c4882

## pr

https://github.com/dphaener/rhiza/pull/506

## result

14 checks passed: changes SUCCESS, go-install SUCCESS, snapshot SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #506 (https://github.com/dphaener/rhiza/pull/506) passed CI at 1dad0c4

## status

verified

## Links

### Out

- verifies [auto-release](../capability/rhiza/distribution/auto-release.md)
- verifies [major-release-signoff](../capability/rhiza/distribution/major-release-signoff.md)
- verifies [distribution](../capability/rhiza/distribution.md)
- verifies [release-goreleaser](../capability/rhiza/distribution/release-goreleaser.md)
- verifies [release-provenance](../capability/rhiza/distribution/release-provenance.md)
- verifies [release-channel-manifest](../capability/rhiza/distribution/release-channel-manifest.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
