# version-channel-check-qa-2

evidence · `01M77KT6NQYCWFESYNE6D31V5J` · actor [claude-code-32b9d78a](../actor/claude-code-32b9d78a.md)

## result

PASS

## statement

QA at b380a501f (cycle 2): built rhiza with go build -trimpath and ran `rhiza version --channel` against a local python http.server channel (RHIZA_TEST_CHANNEL_BASE_URL) with RHIZA_TEST_VERSION set. The cycle-1 ordering bug is fixed. These report newer: v0.2.0-5-gabc1234 over v0.2.0; v0.2.1-12-gdef over v0.2.1-5-gabc; v0.2.1-main.10+abc over main.9; v0.2.0-rc.10 over rc.2; v0.2.0 over rc.10; v0.10.0 over v0.9.0; v0.2.1-main.3 over v0.2.0-7-gabc-dirty; v0.2.0 over dev and over the Go pseudo-version build. These report not newer: each reverse pair, the same version, v0.2.0 against v0.3.0, rc.1 against v0.2.0, and main.3+abc against main.3+abc-dirty. --json gives channel/current/newest/newer. An unknown channel exits 2 with a usage error. A malformed manifest and an unreachable base (127.0.0.1:1) exit 1 with a clean error. Nothing was written in the empty cwd. Against the real published channel, release reports v0.1.0 as not newer than the v0.1.1 pseudo-version build. go test -trimpath ./cmd/rhiza -run TestScripts/version-channel passes.

## Links

### Out

- verifies [version-channel-check](../capability/rhiza/distribution/version-channel-check.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
