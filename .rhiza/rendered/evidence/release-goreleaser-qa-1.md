# release-goreleaser-qa-1

evidence · `01M6Z7CM84AQ7334BCCSQ3WCY1` · actor [claude-code-9da7d8ef](../actor/claude-code-9da7d8ef.md)

## result

PASS

## statement

QA at 527292df with GoReleaser v2.20.0 in a scratch worktree: goreleaser check exited 0; goreleaser release --snapshot --clean exited 0 and wrote rhiza_0.0.0-SNAPSHOT-527292df_{linux,darwin}_{amd64,arm64}.tar.gz plus checksums.txt (4 sha256 lines), and the workflow's four-archive ls loop passed. In a separate clone tagged v0.1.0, goreleaser release --clean --skip=publish,validate exited 0 and wrote rhiza_0.1.0_* for all four targets plus checksums.txt; the extracted linux_amd64 binary is static, stripped and runs (rhiza version, rhiza --help). YAML parse: release.yml triggers only on push (tags v*) with contents: write; release-snapshot.yml and go-install.yml trigger on pull_request only. Publishing to a real GitHub release was not exercised.

## Links

### Out

- verifies [release-goreleaser](../capability/rhiza/distribution/release-goreleaser.md)

## History

- `01M6W5C84ZKXW2C2AE3YWKZ1XD` Release pipeline: GoReleaser builds linux and darwin binaries for amd64 and arm64 on v* tags and publishes them as a GitHub release, rhiza version reports the tag, and go install of the module path works
