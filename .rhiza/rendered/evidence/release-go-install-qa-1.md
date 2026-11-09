# release-go-install-qa-1

evidence · `01M6Z7CMF9Y3K9ZXC32F47WG2S` · actor [claude-code-9da7d8ef](../actor/claude-code-9da7d8ef.md)

## result

PASS

## statement

QA at 527292df: go.mod's module is github.com/dphaener/rhiza and the CI replace grep finds no replace directive. go install github.com/dphaener/rhiza/cmd/rhiza@v0.1.0 with a fresh GOBIN and an empty GOMODCACHE, GOPROXY=file://<local proxy holding only this module>,https://proxy.golang.org, GOSUMDB off, exited 0: every dependency came from the public proxy, and the binary printed 'rhiza v0.1.0'. The go-install.yml step run by hand (fresh GOBIN, go install -mod=readonly ./cmd/rhiza, then rhiza version) exited 0.

## Links

### Out

- verifies [release-go-install](../capability/rhiza/distribution/release-go-install.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
