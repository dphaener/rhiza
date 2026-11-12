# m33-release-interview

source · `01M76AY599871MGQ6E7XCVAK6X`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Let's chat about this one. We should have baked in release requirements. What triggers a release, etc.",
    "where": "on how an official release is cut"
  },
  {
    "by": "owner",
    "text": "#1, and when does a major release happen?",
    "where": "on what triggers a release; #1 is the agent's option 'every finished milestone is a minor release, fix-only patch releases in between'"
  },
  {
    "by": "owner",
    "text": "Fully automatic for minor and patch, owner gated for major releases",
    "where": "on who tags a release once the gates pass"
  },
  {
    "by": "owner",
    "text": "If we're sticking to semver, then a breaking change is always a major release. Does semver allow for a major release that is NOT a breaking change?",
    "where": "on when a major release happens"
  },
  {
    "by": "owner",
    "text": "Batched, unless the bugfix is marked as severe and needs to be released immediately.",
    "where": "on when an automatic patch release is cut"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Rolling prerelease: one GitHub prerelease with a fixed tag (e.g. 'latest') that each merge to main replaces",
    "where": "latest channel home"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Strict from v1: from v1.0.0 any release carrying a breaking change is a major, owner-gated; the owner may also declare a major without one. Before v1 breaking changes bump the minor, and v1.0.0 is cut when the owner says so.",
    "where": "breaking changes and majors"
  }
]
```

## raw_sha256

dac0ec51d5c1542176f1e4fd253a6a5dc2a6e4dd68d041e9564bc0fdf59f9caa

## session

41425f42-c4ee-4701-bd80-01a1cc544960

## summary

Planned milestone 33, release channels for Rhiza binaries (milestone 32's auto-update will rely on it). Channels: release (official v* GitHub releases) and latest (one rolling GitHub prerelease tagged 'latest', replaced on every merge to main); main stays the local build from main. A latest build reports the next patch as a semver prerelease, e.g. v0.1.1-main.42+abc1234 after v0.1.0 (the agent first offered a base that sorted before the release and corrected it). Clients find the newest build through a JSON manifest asset per channel at a fixed download URL, and verify the archive's sha256 against it; CI also publishes GitHub build-provenance attestations. Client side here: a resolver package and a read-only `rhiza version --channel`, plus an install.sh; milestone 32 owns the daemon's update mechanics. The owner asked for baked-in release requirements: every finished milestone is a minor release; bug fixes are batched into a patch at most once per 24h soak, unless a fix is marked severe, which releases as soon as it lands (skipping the soak, never CI). Gates: CI green on the commit, 24h soak of that commit's latest build on the VM daemon with no bug filed against it, the milestone closed, notes generated from the model. Minor and patch releases are fully automatic; majors are owner-gated. Semver is strict from v1.0.0 (a breaking change is a major; the owner may also declare one); before v1 breaking changes bump the minor and the owner decides when v1.0.0 is cut. A change is breaking when a change set is marked so at plan or review, or when the release check detects a breaking diff in API schemas, settings or log format. v0.1.0 is cut by hand once the release manifest lands. All eight change sets are planned now.

## title

Milestone 33 plan: release channels and release policy, with the owner

## Links

### In

- derived-from [m33-release-requirements](../decision/m33-release-requirements.md)
- derived-from [m33-milestone-is-minor](../decision/m33-milestone-is-minor.md)
- derived-from [m33-auto-minor-patch](../decision/m33-auto-minor-patch.md)
- derived-from [m33-breaking-is-major](../decision/m33-breaking-is-major.md)
- derived-from [m33-batched-patches](../decision/m33-batched-patches.md)
- derived-from [answer-41425f42-1](../decision/answer-41425f42-1.md)
- derived-from [answer-41425f42-2](../decision/answer-41425f42-2.md)
- derived-from [answer-41425f42-3](../decision/answer-41425f42-3.md)
- derived-from [answer-41425f42-4](../decision/answer-41425f42-4.md)
- derived-from [answer-41425f42-5](../decision/answer-41425f42-5.md)
- derived-from [answer-41425f42-7](../decision/answer-41425f42-7.md)
- derived-from [answer-41425f42-8](../decision/answer-41425f42-8.md)
- derived-from [answer-41425f42-11](../decision/answer-41425f42-11.md)
- derived-from [answer-41425f42-13](../decision/answer-41425f42-13.md)
- derived-from [answer-41425f42-15](../decision/answer-41425f42-15.md)
- derived-from [answer-41425f42-16](../decision/answer-41425f42-16.md)
- derived-from [answer-41425f42-17](../decision/answer-41425f42-17.md)
- derived-from [answer-41425f42-18](../decision/answer-41425f42-18.md)
- derived-from [answer-41425f42-19](../decision/answer-41425f42-19.md)
- derived-from [answer-41425f42-20](../decision/answer-41425f42-20.md)

## History

- `01M765854D7DGJDBW5CYGJ4SR1` Release channels for Rhiza binaries. Owner, 2026-11-11: "As a user, I would like to be able to download a Rhiza binary from a stable release channel. Milestone 32 is planning some auto update features as well and will rely on this. The daemon will have a setting that can auto update on the following settings: off, release, latest, and main. Off is self explanatory. Release is the latest official release. Latest is the latest build (All merges to main will create a build that is available for download, but will not necessarily be a release). Main is what we are doing here right now, building directly from main locally and installing from that, the most bleeding edge of bleeding edge. Milestone 32 will handle all of the mechanics of auto updating, we just need to make all the release channels happen." Orchestrator's findings to start from (verify, don't trust): milestone 30 shipped a release pipeline (#407): .goreleaser.yaml, .github/workflows/release.yml (publishes a GitHub release on v* tags only), release-snapshot.yml (dry-runs the archives on PRs that touch the release config), go-install.yml; four archives (Linux and macOS, amd64 and arm64). No release has been published yet (gh release list is empty, no v* tags), the README points to the releases page and go install ...@latest. The daemon's self-update today (rhiza.daemon.selfUpdate, internal/api/daemon.go:320-345, internal/daemon SelfUpdate) fetches main and runs make install in the main checkout. Milestone 32 (workflow cleanup and language-neutral defaults, request 01M75W4E, plan run 01M75W8Y, interview in progress) renames it rhiza.daemon.autoUpdate with values off and main, leaves release and latest to this milestone, and will depend on this milestone. In scope here: the channels themselves (official releases and a latest build per merge to main, where they are published, how a client finds the newest of each and verifies it), not the daemon's update mechanics.
