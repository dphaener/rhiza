# auto-release

capability · `01M76B7EAW7JX0E9RNRBMXM97G` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... -run TestCut with a fake git tags and pushes once when every gate passes, never for a major, never twice for one version and never when a gate fails; go test -trimpath ./tools/release/... -run TestRecordDaemonRun with a fake daemon status appends a (commit, started) entry to the daemon-run history only when the commit is new; release.yml uses the tag annotation as the notes; the timer unit is in the repo and the VM runbook names it

## statement

A systemd user timer on the VM runs tools/release hourly in the root main checkout; each run first records the root daemon's running commit and start time in the daemon-run history the soak gate reads, and when check passes for a minor or patch it pushes an annotated tag carrying the notes, and the release job publishes them as the release body

## Links

### Out

- decided-by [m33-auto-minor-patch](../../../decision/m33-auto-minor-patch.md)

### In

- relates-to [major-signoff-lightweight-tag](../../../question/major-signoff-lightweight-tag.md) — a signed-off major may never get its annotated tag
- verifies [auto-release-verified](../../../evidence/auto-release-verified.md)
- verifies [auto-release-qa-1](../../../evidence/auto-release-qa-1.md)
- verifies [auto-release-verified-2](../../../evidence/auto-release-verified-2.md)
- verifies [auto-release-qa-2](../../../evidence/auto-release-qa-2.md)
- verifies [pr-506-automatic-releases](../../../evidence/pr-506-automatic-releases.md)

## History

- `01M76B0M0A1T1SN6KHMG5JVDVR` Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release
