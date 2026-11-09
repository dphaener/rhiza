# release-version-stamp-qa-1

evidence · `01M6Z7CMB1NY6AHH4TDZVFW2XB` · actor [claude-code-9da7d8ef](../actor/claude-code-9da7d8ef.md)

## result

PASS

## statement

QA at 527292df: the GoReleaser build from a clone tagged v0.1.0 prints 'rhiza v0.1.0' (ldflags -X main.version={{ .Tag }}); go install github.com/dphaener/rhiza/cmd/rhiza@v0.1.0 through a local file GOPROXY serving a v0.1.0 module zip (git archive of the tag) prints 'rhiza v0.1.0' from the module version, and go version -m shows mod github.com/dphaener/rhiza v0.1.0. go test -trimpath ./cmd/rhiza -run TestVersion -count=1 passed (TestVersion, TestVersionGoReleaserStamp).

## Links

### Out

- verifies [release-version-stamp](../capability/rhiza/distribution/release-version-stamp.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
