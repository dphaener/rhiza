# release-version-stamp-verified

evidence · `01M6Z6JMM449V6TPEMAJ42FDP3` · actor [claude-code-e893d2fc](../actor/claude-code-e893d2fc.md)

## result

PASS

## statement

At 7310d511, go test -trimpath ./cmd/rhiza -run TestVersion -v -count=1 passed: TestVersion's table covers the ldflags path (v0.1.0 wins, including over a module v0.2.0), the go install module-version path (v0.1.0) and the dev fallback for (devel) and empty; TestVersionGoReleaserStamp passed. .goreleaser.yaml line 18 sets -s -w -X main.version={{ .Tag }}. A build with -ldflags "-X main.version=v0.1.0" prints 'rhiza v0.1.0'. gofmt and go vet on cmd/rhiza are clean.

## Links

### Out

- verifies [release-version-stamp](../capability/rhiza/distribution/release-version-stamp.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
