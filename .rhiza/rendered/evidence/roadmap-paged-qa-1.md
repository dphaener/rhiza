# roadmap-paged-qa-1

evidence · `01M6J17XEPDH05Z04N7SM1ZDM8` · actor [claude-code-17e5e990](../actor/claude-code-17e5e990.md)

## result

PASS

## statement

At b63f7c18 with a scratch build on this repo (4 groups, 18 change sets): paging roadmap with --limit 1, 4, 5 and 7 to the end gave exactly the (group kind, intent, change set) entries of --all, in order, each once. --limit 3 split milestone-25 and page 2 repeated its '## milestone-25' heading before Backlog. --limit 0 equals the default; --limit -1, --limit 1001 and --cursor zz exited 2 ('bad cursor "zz": give a next_cursor from roadmap'). roadmap --json and api call roadmap printed the same bytes. TestRoadmapOverAPageIsReadWhole (health), TestReadsARoadmapOverAPageWhole (roadmap app), TestRoadmapTabReadsARoadmapOverAPageWhole (model app) and roadmap-paged.txtar pass.

## Links

### Out

- verifies [roadmap-paged](../capability/rhiza/views/paged-lists/roadmap-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
