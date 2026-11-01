# inbox-overlap-items-qa-1

evidence · `01M6A9VHQQ6PC4X525YZ53GDKA` · actor [claude-code-687f5612](../actor/claude-code-687f5612.md)

## result

PASS

## statement

QA at 158cf26b with a scratch build in a fresh git repo, RHIZA_* env cleared. With two proposed sets on thing's statement, rhiza inbox lists '5a74731 overlap change sets A and B', 'target: thing', 'both write field:statement' and the next slice command; the ID was stable across calls. inbox --json has kind overlap, change A, change_title, other B, target thing, detail and next, and no run; rhiza api call inbox prints the same bytes. rhiza inbox open 5a74731 exits 2 with 'is an overlap between change sets A and B: it has no tab to open' and the next command. After A landed the inbox read 'Nothing waits on you'. In a Rex pane, rhiza tui app inbox showed an Overlaps group with 'Other edit · with <B> · other' / 'both write field:statement' and 'Retire goal · with <B> · goal' / 'both write linked, retired'; pressing o showed 'open failed: item 24dc55b is an overlap ...' and the app kept running. The daemon's notice text was not run against a live daemon (it would start runs for accepted sets); it is covered by TestInboxNoticeOverlap.

## Links

### Out

- verifies [inbox-overlap-items](../capability/rhiza/orchestration/inbox/inbox-overlap-items.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
