# root-daemon-self-update

capability · `01M6RN85F2T3NMVDRKV706VSZB` · parent [daemon](../daemon.md)

## acceptance

internal/daemon selfupdate tests: with self_update_project set, the fake git and make run in that project's checkout; unset, neither runs.

## statement

With selfUpdate on, the root daemon keeps itself current from the project rhiza.toml names as self_update_project: it fast-forwards that project's main checkout and runs make update there when origin/main moves, as today's self-update does; with none named it never self-updates.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)
- supersedes daemon-self-update (retired)

### In

- verifies [root-daemon-self-update-verified](../../../../evidence/root-daemon-self-update-verified.md)
- verifies [root-daemon-self-update-qa-2](../../../../evidence/root-daemon-self-update-qa-2.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
