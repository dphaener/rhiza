# testscripts-run-under-root

capability · `01M6RNXQCZZ7E21TJ70B4RCRWW` · parent [runtime](../runtime.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestScripts/^(run-start|run-start-queue|daemon-settings|daemon-pause|daemon-event-hook|stage-autonomy)$' passes with runs-need-root and root-daemon-lock in force and those scripts unedited (daemon-repo-session is rewritten by root-daemon-rex-session's leaf); a script running rhiza root show in $WORK/repo sees project repo, and in $NOROOT gets exit 2 naming rhiza root init.

## statement

The testscript Setup in cmd/rhiza/main_test.go puts every script under a Rhiza root: it makes each script's $WORK its own root, writing $WORK/rhiza.toml that registers a project named repo at $WORK/repo directly (not through rhiza project add), so the existing scripts that start runs, queue requests or run the daemon in the $WORK/repo repo keep passing unedited, and it gives scripts a folder outside any root ($NOROOT) for the outside-a-root cases.

## Links

### Out

- decided-by [m29-root-required](../../../decision/m29-root-required.md)
- decided-by [m29-test-root-per-script](../../../decision/m29-test-root-per-script.md)
- decided-by [m29-rewrite-daemon-repo-session](../../../decision/m29-rewrite-daemon-repo-session.md)

### In

- verifies [testscripts-run-under-root-verified](../../../evidence/testscripts-run-under-root-verified.md)
- verifies [testscripts-run-under-root-verified-2](../../../evidence/testscripts-run-under-root-verified-2.md)
- verifies [testscripts-run-under-root-qa-2](../../../evidence/testscripts-run-under-root-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
