# changes-app-verified

evidence · `01M6BDVQ8N953XYQGZZ760RD6N` · actor [claude-code-0a11f7de](../actor/claude-code-0a11f7de.md)

## result

PASS

## statement

Verifier at ce64162f: go test ./cmd/rhiza -run 'TestScripts/(tui-changes-app|rex-install)$' -count=1 passed. tui-changes-app.txtar shows --once listing the unseen landed change sets newest first with "pr": 7, a second --once printing "change_sets": [], and a later-landed gamma appearing alone, then none again, with no socket or daemon lock. go test ./internal/tui/changesapp passed its three fake-reader unit tests: newest-first order across pages, seen IDs not listed again, a later change set listed alone, and the seen file holding sorted IDs. rex-install.txtar greps rhiza_app_changes on ["changes"] = "g". I also built rhiza from HEAD~2 and from ce64162f and ran rex install with a fake rex; diffing the generated key maps showed only the added ["changes"] = "g" line, so no other key moved. Reading the diff: launcher.App gains a Key that rexDefaultKeys reserves before auto-picking keys, and seen IDs go to .git/rhiza/feed.json under git rev-parse --git-common-dir. make check exited 0.

## Links

### Out

- verifies changes-app (retired)

## History

- `01M6B3NCDEX88PC3NFRW16YTVZ` rhiza feed lists what landed in land order with PR numbers, and a Rex changes app shows what landed since the owner last looked
