# pr-407-release-pipeline

evidence · `01M6Z7R6SWJ821RREWF4BQ8H33` · actor [dphaener](../actor/dphaener.md)

## head

7849b5c92ef7754a5e3d5f9b2617714d10b6ccf5

## pr

https://github.com/dphaener/rhiza/pull/407

## result

12 checks passed: lint SUCCESS, go-install SUCCESS, snapshot SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #407 (https://github.com/dphaener/rhiza/pull/407) passed CI at 7849b5c

## status

verified

## Links

### Out

- verifies [release-goreleaser](../capability/rhiza/distribution/release-goreleaser.md)
- verifies [release-version-stamp](../capability/rhiza/distribution/release-version-stamp.md)
- verifies [release-go-install](../capability/rhiza/distribution/release-go-install.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
