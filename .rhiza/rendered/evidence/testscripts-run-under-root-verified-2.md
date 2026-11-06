# testscripts-run-under-root-verified-2

evidence · `01M6S7JX515WS8RBYM14RH9896` · actor [claude-code-99bb450f](../actor/claude-code-99bb450f.md)

## result

PASS

## statement

At 30479135, go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^(run-start|run-start-queue|daemon-settings|daemon-pause|daemon-event-hook|stage-autonomy|testscript-root)$' passed all seven. git log main..HEAD shows none of the six listed scripts changed. runs-need-root is in force: requireProject is at internal/api/run.go:534 and request.go:103. root-daemon-lock belongs to another change set and is not built at this commit, so only runs-need-root could be checked as in force. Setup in cmd/rhiza/main_test.go writes $WORK/rhiza.toml ([projects.repo] path="repo") directly and sets $NOROOT to a temp dir outside any root. testscript-root.txtar asserts that root show in $WORK/repo prints project: repo, and that root show in $NOROOT fails with 'not inside a Rhiza root' and next: rhiza root init. A scratch build run in a folder outside any root exited with code 2 and the same message.

## Links

### Out

- verifies [testscripts-run-under-root](../capability/rhiza/runtime/testscripts-run-under-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
