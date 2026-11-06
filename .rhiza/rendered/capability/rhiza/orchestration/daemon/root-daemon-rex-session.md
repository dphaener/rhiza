# root-daemon-rex-session

capability · `01M6RN85G5TY8WJ23TXEGTZFGV` · parent [daemon](../daemon.md)

## acceptance

A testscript runs rhiza daemon --once on a root with two projects, one added with --name other than its folder, against a fake Rex with REX_SESSION unset: it creates rhiza-<name> once for each project missing from rex session ls, each run driver it starts carries REX_SESSION=rhiza-<registered name>, with the sessions present it creates none, and rhiza run open for the renamed project's run targets rhiza-<registered name>.

## statement

Every Rex tab the root daemon's runs open (agent tabs; plan-request runs and run drivers start headless) opens in its project's fixed Rex session, rhiza-<project>, where <project> is the name the project is registered under in rhiza.toml (its folder, or project add --name), not the checkout's folder; the daemon creates each session on its Rex server when missing and passes it as REX_SESSION to every run driver it starts for that project, and rhiza inbox open and run open use the same registered name to find a run's panes.

## Links

### Out

- decided-by [m29-rex-session-per-project](../../../../decision/m29-rex-session-per-project.md)
- supersedes daemon-repo-session (retired)

### In

- verifies [root-daemon-rex-session-verified](../../../../evidence/root-daemon-rex-session-verified.md)
- verifies [root-daemon-rex-session-qa-1](../../../../evidence/root-daemon-rex-session-qa-1.md)
- verifies [root-daemon-rex-session-qa-2](../../../../evidence/root-daemon-rex-session-qa-2.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
