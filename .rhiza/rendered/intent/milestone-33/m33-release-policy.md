# m33-release-policy

intent · `01M76AYSVSA60HBQ48YZFSVRS3` · parent [milestone-33](../milestone-33.md)

## statement

Releases are cut by rule: a finished milestone is a minor, landed fixes a batched patch (a severe fix at once), every release passes CI, a 24h soak and the milestone-closed gate with notes from the model; minor and patch are automatic, a major is owner-gated and strict semver holds from v1.0.0.

## Links

### Out

- decided-by [m33-release-requirements](../../decision/m33-release-requirements.md)
- decided-by [m33-milestone-is-minor](../../decision/m33-milestone-is-minor.md)
- decided-by [m33-auto-minor-patch](../../decision/m33-auto-minor-patch.md)
- decided-by [m33-breaking-is-major](../../decision/m33-breaking-is-major.md)
- decided-by [m33-batched-patches](../../decision/m33-batched-patches.md)

### In

- serves change set `01M76B077D23CJ28C690Y166VY` (Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against)
- serves change set `01M76B0BSVESMEFD39RTW517R3` (Breaking-change detection: compare the API schemas (describe --json), daemon settings and log format at HEAD with the last release tag and report each breaking difference)
- serves change set `01M76B0FZA3N7ZHAYTBN76MVSV` (rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model)
- serves change set `01M76B0M0A1T1SN6KHMG5JVDVR` (Automatic releases: the daemon cuts a minor or patch release when rhiza release check passes (tags it, the release workflow publishes it with the generated notes), and raises an owner sign-off in the inbox for a major release)

## History

- `01M765854D7DGJDBW5CYGJ4SR1` Release channels for Rhiza binaries. Owner, 2026-11-11: "As a user, I would like to be able to download a Rhiza binary from a stable release channel. Milestone 32 is planning some auto update features as well and will rely on this. The daemon will have a setting that can auto update on the following settings: off, release, latest, and main. Off is self explanatory. Release is the latest official release. Latest is the latest build (All merges to main will create a build that is available for download, but will not necessarily be a release). Main is what we are doing here right now, building directly from main locally and installing from that, the most bleeding edge of bleeding edge. Milestone 32 will handle all of the mechanics of auto updating, we just need to make all the release channels happen." Orchestrator's findings to start from (verify, don't trust): milestone 30 shipped a release pipeline (#407): .goreleaser.yaml, .github/workflows/release.yml (publishes a GitHub release on v* tags only), release-snapshot.yml (dry-runs the archives on PRs that touch the release config), go-install.yml; four archives (Linux and macOS, amd64 and arm64). No release has been published yet (gh release list is empty, no v* tags), the README points to the releases page and go install ...@latest. The daemon's self-update today (rhiza.daemon.selfUpdate, internal/api/daemon.go:320-345, internal/daemon SelfUpdate) fetches main and runs make install in the main checkout. Milestone 32 (workflow cleanup and language-neutral defaults, request 01M75W4E, plan run 01M75W8Y, interview in progress) renames it rhiza.daemon.autoUpdate with values off and main, leaves release and latest to this milestone, and will depend on this milestone. In scope here: the channels themselves (official releases and a latest build per merge to main, where they are published, how a client finds the newest of each and verifies it), not the daemon's update mechanics.
