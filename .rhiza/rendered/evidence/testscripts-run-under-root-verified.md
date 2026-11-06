# testscripts-run-under-root-verified

evidence · `01M6S2ZZ688KHE398WAH6TGFZZ` · actor [claude-code-66235903](../actor/claude-code-66235903.md)

## result

PASS

## statement

At 1648261d, go test -trimpath ./cmd/rhiza -run 'TestScripts/^(run-start|run-start-queue|daemon-settings|daemon-pause|daemon-event-hook|daemon-repo-session|stage-autonomy|testscript-root|runs-need-root)$' passed every script. git diff 85470cb5^..1648261d shows those seven scripts unedited, while runs-need-root (85470cb5) is in force. The Setup in cmd/rhiza/main_test.go writes $WORK/rhiza.toml registering project repo at $WORK/repo and sets $NOROOT to a temp folder outside any root, per decision m29-test-root-per-script. In the new testscript-root.txtar, rhiza root show in $WORK/repo prints root: and project: repo; in $NOROOT it exits non-zero with 'not inside a Rhiza root' and next: rhiza root init. root-daemon-lock belongs to change set 01M6RN3RCGAZ9RNTWTZPKFM61M and is not built at this commit, so the daemon scripts were checked under the root without it.

## Links

### Out

- verifies [testscripts-run-under-root](../capability/rhiza/runtime/testscripts-run-under-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
