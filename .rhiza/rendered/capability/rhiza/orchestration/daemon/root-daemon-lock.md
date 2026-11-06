# root-daemon-lock

capability · `01M6RN858XDJS5W8FZ3QZ8P2PF` · parent [daemon](../daemon.md)

## acceptance

internal/daemon lock tests: a second lock on the same root fails, from a project dir and from the root alike; a testscript gets exit 2 with the next from rhiza daemon in $NOROOT, outside any root; the existing daemon scripts (daemon-settings, daemon-pause, daemon-event-hook, daemon-repo-session, stage-autonomy) pass unedited under the root testscripts-run-under-root gives them.

## statement

rhiza daemon runs on the root, never on a single repo: run anywhere under the root it resolves the root and holds one lock in the root's state folder (<root>/.rhiza-root/), so a second daemon on the same root refuses to start, and outside a root it exits 2 with next rhiza root init.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)
- decided-by [m29-root-required](../../../../decision/m29-root-required.md)
- supersedes daemon-lock (retired)

### In

- verifies [root-daemon-lock-verified](../../../../evidence/root-daemon-lock-verified.md)
- verifies [root-daemon-lock-qa-1](../../../../evidence/root-daemon-lock-qa-1.md)
- verifies [root-daemon-lock-qa-2](../../../../evidence/root-daemon-lock-qa-2.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
