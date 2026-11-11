# m33-latest-channel

intent · `01M76AYSSQZ8JHN8SNH038CERR` · parent [milestone-33](../milestone-33.md)

## statement

Every merge to main replaces one rolling 'latest' GitHub prerelease with that commit's four archives, checksums and manifest, versioned as the next patch's semver prerelease.

## Links

### In

- serves change set `01M76AZHGY605TXF68M27W4SZV` (Latest channel: every merge to main builds the four archives as the next patch semver prerelease (vX.Y.Z+1-main.N+sha) and replaces one rolling latest GitHub prerelease with them, their checksums, manifest and attestations)

## History

- `01M765854D7DGJDBW5CYGJ4SR1` Release channels for Rhiza binaries. Owner, 2026-11-11: "As a user, I would like to be able to download a Rhiza binary from a stable release channel. Milestone 32 is planning some auto update features as well and will rely on this. The daemon will have a setting that can auto update on the following settings: off, release, latest, and main. Off is self explanatory. Release is the latest official release. Latest is the latest build (All merges to main will create a build that is available for download, but will not necessarily be a release). Main is what we are doing here right now, building directly from main locally and installing from that, the most bleeding edge of bleeding edge. Milestone 32 will handle all of the mechanics of auto updating, we just need to make all the release channels happen." Orchestrator's findings to start from (verify, don't trust): milestone 30 shipped a release pipeline (#407): .goreleaser.yaml, .github/workflows/release.yml (publishes a GitHub release on v* tags only), release-snapshot.yml (dry-runs the archives on PRs that touch the release config), go-install.yml; four archives (Linux and macOS, amd64 and arm64). No release has been published yet (gh release list is empty, no v* tags), the README points to the releases page and go install ...@latest. The daemon's self-update today (rhiza.daemon.selfUpdate, internal/api/daemon.go:320-345, internal/daemon SelfUpdate) fetches main and runs make install in the main checkout. Milestone 32 (workflow cleanup and language-neutral defaults, request 01M75W4E, plan run 01M75W8Y, interview in progress) renames it rhiza.daemon.autoUpdate with values off and main, leaves release and latest to this milestone, and will depend on this milestone. In scope here: the channels themselves (official releases and a latest build per merge to main, where they are published, how a client finds the newest of each and verifies it), not the daemon's update mechanics.
