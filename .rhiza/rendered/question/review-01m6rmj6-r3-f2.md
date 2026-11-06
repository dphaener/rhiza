# review-01m6rmj6-r3-f2

question · `01M6RP4X21TN1QGFWBWR8P2WJZ`

## statement

Plan review round 3 left finding f2 (contradiction on testscripts-run-under-root): Its acceptance, and root-daemon-lock's, require daemon-repo-session.txtar to pass unedited. That script asserts the Rex session rhiza-repo, named after the checkout's folder. This change set retires daemon-repo-session and replaces it with root-daemon-rex-session, which names the session after the project's registered name (the $WORK folder, script-daemon-repo-session, under this Setup). The script tests a retired capability and would see a different session, so it cannot pass unedited. The critic confirmed it: daemon-repo-session.txtar asserts 'rex session new rhiza-repo' (folder-named). Change set 01M6RN3RCGAZ9RNTWTZPKFM61M retires 01M6FES3PKQCSNG9SEECK4G1R7 as superseded by root-daemon-rex-session (rhiza-<project>, the registered name). Yet testscripts-run-under-root and root-daemon-lock require that script to pass unedited, testing a retired capability. How should it be settled?

## Links

### In

- answers [m29-rewrite-daemon-repo-session](../decision/m29-rewrite-daemon-repo-session.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
