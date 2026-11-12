# latest-build-version-verified

evidence · `01M77HDX2A4SBTT4MSRWYG3JRS` · actor [claude-code-fbfd878e](../actor/claude-code-fbfd878e.md)

## result

PASS

## statement

Verifier at 77cd67581: go test -trimpath -count=1 ./tools/release/... passes (TestLatestVersion, TestLatestVersionOrdering, TestGitLatestVersion); go vet and gofmt clean. TestGitLatestVersion builds a hermetic repo and computes v0.0.1-main.2+sha with no tag, v0.3.0 with the tag at HEAD, and v0.3.1-main.2+sha with the tag 2 commits back. The test's semverLess orders v0.1.0 < v0.1.1-main.3 < main.4 < main.10 < v0.1.1. Cross-checked independently with golang.org/x/mod/semver v0.37.0: all those versions are valid semver and Compare gives the same order (also v0.0.0 < v0.0.1-main.12+abc1234 < v0.0.1). go run ./tools/release version on this checkout printed v0.0.1-main.477+77cd675. Diff read: version.go formats vX.Y.(Z+1)-main.N+<7-char sha> from git describe --match v*.*.* and rev-list --count.

## Links

### Out

- verifies [latest-build-version](../capability/rhiza/distribution/latest-build-version.md)

## History

- `01M76AZHGY605TXF68M27W4SZV` Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations
