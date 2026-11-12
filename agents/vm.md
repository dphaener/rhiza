# Operating the VM

Read this when you operate the VM that runs Rhiza's daemon (the VM `rhiza`,
Ubuntu, 4 cores shared by the daemon and every live run). The orchestrator's
watch loop is in [orchestrator.md](orchestrator.md); this doc holds the
lessons about the machine itself.

## One root

The VM runs one root, `~/code`, under the `rhiza-root.service` systemd user
service. `~/code/rhiza.toml` registers the projects and holds the machine
settings. Project settings (autonomy, pause) stay in each repo's git config;
`rhiza daemon set` writes each key to its own scope. Builds refuse to run
the daemon outside a root. Never start a second daemon for this repo (on a
laptop or from another clone): two daemons would double-start change sets.

## Load

Runs are CPU heavy (Go test suites, browser specs). When the 5-minute load
average stays high, lower `rhiza.daemon.maxRuns`:

```bash
rhiza daemon set rhiza.daemon.maxRuns <n>
```

Measured on the 4 cores: 3 runs hold the load around 10 to 16, 4 runs reach
16 to 25, and 5 runs pass 25, make `rhiza health` take minutes and fail Go
builds with missing `$WORK` files. Look at the 5-minute average, not the
1-minute one: short spikes to 18 to 23 pass. Lowering the cap never stops
runs already live.

## Automatic releases

The systemd user timer `rhiza-release.timer` (units in `script/systemd/`)
runs `rhiza-release.service` hourly in the root main checkout
`~/code/rhiza`. Each run first records the root daemon's running commit and
start time in the daemon-run history (`$RHIZA_RELEASE_RUNS`, default in the
unit `~/.local/state/rhiza/daemon-runs.jsonl`), which the soak gate reads,
adding a line only for a commit not yet in it. Then `tools/release cut` pushes
an annotated tag with the notes when `release check` passes for a minor or
patch; the release job publishes the notes as the release body. A major
release is never tagged: it waits for the owner. Install:

```bash
mkdir -p ~/.config/systemd/user
cp ~/code/rhiza/script/systemd/rhiza-release.* ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable --now rhiza-release.timer
```

Read a run with `journalctl --user -u rhiza-release.service`.

## Pinned Rex package

The Rex package is pinned. Ubuntu's universe ships an unrelated Perl `rex`
(Rexify) with a higher version, and unattended-upgrades once replaced the
Rex terminal with it: every `rex` call then printed "No Rexfile found" and
the daemon could not open agent tabs. `/etc/apt/preferences.d/rex` pins the
Rex package source above Ubuntu's, and `/etc/apt/apt.conf.d/51unattended-no-rex`
keeps unattended-upgrades off it. If a run fails with "No Rexfile found", the
wrong `rex` is on `PATH`: check `apt-cache policy rex` and `dpkg -S /usr/bin/rex`.

## Stray log entries

When a run worktree is reset to main, its untracked `.rhiza/log` entries are
moved aside to `~/orchestrator/stray-logs/<worktree-slug>/.rhiza/log/`. The
model then shows a change set without those writes, and a review re-reports
findings already fixed. If a change set's ops look older than its rounds
suggest, look in `~/orchestrator/stray-logs/` and copy that change set's
entries back; check they chain onto the current head, then confirm with
`rhiza change show`.

## Flakes

Reproduce a timing flake by injecting a delay at the suspected race (wrap the
fetch, timer or job), never by CPU load. Busy loops and parallel stress runs
on the shared host slow every live run, and rarely trigger the flake anyway.
If a run's builder is burning CPU on load generators, stop those processes.
