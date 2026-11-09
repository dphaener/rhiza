# release-go-install-verified

evidence · `01M6Z6QJJ3MVRZ4Z6NH5V62F43` · actor [claude-code-092ae07f](../actor/claude-code-092ae07f.md)

## result

PASS

## statement

At 397aa710 read the diff: new .github/workflows/go-install.yml runs on PRs to main touching go.mod/go.sum/cmd/internal; it fails if go.mod has a replace directive (grep '^[[:space:]]*replace[[:space:]]', which also catches a 'replace (' block, checked), then sets GOBIN=$RUNNER_TEMP/gobin, rm -rf it, runs go install -mod=readonly ./cmd/rhiza and "$GOBIN/rhiza" version. Workflow YAML parses; runner and action versions match ci.yml. go.mod module is github.com/dphaener/rhiza (matches origin) and has no replace directives. Ran the same steps locally with a fresh GOBIN in a scratch dir: the replace check passed, go install -mod=readonly succeeded and rhiza version printed a version, exit 0.

## Links

### Out

- verifies [release-go-install](../capability/rhiza/distribution/release-go-install.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
