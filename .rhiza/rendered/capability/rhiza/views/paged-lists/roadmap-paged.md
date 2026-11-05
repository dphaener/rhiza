# roadmap-paged

capability · `01M6H3MCPWZBAV528M4JYX6AA8` · parent [paged-lists](../paged-lists.md)

## acceptance

A test pages a roadmap whose group is split by --limit and sees the group's kind and intent repeated on the next page, with every (group, change set) entry appearing exactly once across the pages, so a change set that serves two milestones appears once in each of its two groups; --all matches the concatenated pages; a bad cursor and an out-of-range limit exit 2; a roadmap of more than 100 change sets is read whole by health's roadmap check; roadmap --json and api call roadmap print the same bytes; describe-golden.txtar is regenerated

## statement

rhiza roadmap pages its change sets in group order (milestones, then backlog, then umbrellas; depends-on order within a group) with --limit, --cursor and --all; a page holds only the groups with change sets on it, a group split across pages repeats its kind and intent, and health's roadmap check reads it with all

## Links

### Out

- decided-by [decision-paging-log-roadmap-why](../../../../decision/decision-paging-log-roadmap-why.md)

### In

- verifies [roadmap-paged-verified](../../../../evidence/roadmap-paged-verified.md)
- verifies [roadmap-paged-qa-1](../../../../evidence/roadmap-paged-qa-1.md)
- verifies [pr-228-log](../../../../evidence/pr-228-log.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
