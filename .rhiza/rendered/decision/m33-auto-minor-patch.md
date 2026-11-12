# m33-auto-minor-patch

decision · `01M76AYSXNBZYZEDCBQZF1TXZF`

## question

Once the gates pass, who tags the release?

## quote

Fully automatic for minor and patch, owner gated for major releases

## statement

Minor and patch releases are cut fully automatically once the gates pass; a major release waits for the owner's sign-off

## Links

### Out

- derived-from [m33-release-interview](../source/m33-release-interview.md)

### In

- decided-by [m33-release-policy](../intent/milestone-33/m33-release-policy.md)

## History

- `01M765854D7DGJDBW5CYGJ4SR1` Release channels for Rhiza binaries. Owner, 2026-11-11: "As a user, I would like to be able to download a Rhiza binary from a stable release channel. Milestone 32 is planning some auto update features as well and will rely on this. The daemon will have a setting that can auto update on the following settings: off, release, latest, and main. Off is self explanatory. Release is the latest official release. Latest is the latest build (All merges to main will create a build that is available for download, but will not necessarily be a release). Main is what we are doing here right now, building directly from main locally and installing from that, the most bleeding edge of bleeding edge. Milestone 32 will handle all of the mechanics of auto updating, we just need to make all the release channels happen." Orchestrator's findings to start from (verify, don't trust): milestone 30 shipped a release pipeline (#407): .goreleaser.yaml, .github/workflows/release.yml (publishes a GitHub release on v* tags only), release-snapshot.yml (dry-runs the archives on PRs that touch the release config), go-install.yml; four archives (Linux and macOS, amd64 and arm64). No release has been published yet (gh release list is empty, no v* tags), the README points to the releases page and go install ...@latest. The daemon's self-update today (rhiza.daemon.selfUpdate, internal/api/daemon.go:320-345, internal/daemon SelfUpdate) fetches main and runs make install in the main checkout. Milestone 32 (workflow cleanup and language-neutral defaults, request 01M75W4E, plan run 01M75W8Y, interview in progress) renames it rhiza.daemon.autoUpdate with values off and main, leaves release and latest to this milestone, and will depend on this milestone. In scope here: the channels themselves (official releases and a latest build per merge to main, where they are published, how a client finds the newest of each and verifies it), not the daemon's update mechanics.
