# version-channel-check-qa-1

evidence · `01M77K7Z3VYRT5AMPK9AEHWN5W` · actor [claude-code-a983b6c0](../actor/claude-code-a983b6c0.md)

## result

PASS

## statement

QA at 04a4369f0 with a built binary (go build -trimpath) against a local python http.server channel fixture: release v0.2.0 is reported newer for running v0.1.0 and for v0.2.0-rc.1, and not newer for v0.2.0 and v0.3.0. The --json output is {channel,current,newer,newest}. Unknown channel and misspelled flag exit 2 with next; a 404 or malformed manifest exits 1 with a clean error and --json error. Against the real network, --channel release reports v0.1.0. --channel latest returns 404 (nothing is published there yet). Nothing is written to the cwd. The version-channel and describe-golden testscripts pass. Comparison of git-describe versions on latest is wrong; see the qa finding.

## Links

### Out

- verifies [version-channel-check](../capability/rhiza/distribution/version-channel-check.md)

## History

- `01M76AZSNW1912Y5RD6J4TAHJQ` Channel resolver: a Go package that fetches a channel manifest, picks this platform archive, downloads it and verifies its sha256, and a read-only rhiza version --channel release|latest that shows the newest version there
