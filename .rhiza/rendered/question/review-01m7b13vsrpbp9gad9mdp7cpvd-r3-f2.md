# review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f2

question · `01M7B2GQTEEV1N9NYZVFPZWFGX`

## statement

Plan review round 3 left finding f2 (contradiction on web-project-needs-you): The count includes 'holds' and is to be filled 'from that project's model and runs' and 'read from that project alone' (with a per-project fixture holding 'one hold'), but in the code holds are the daemon's machine-wide items, not a project's: internal/api/inbox.go's InboxItem.Project is 'absent for the daemon's machine-wide items (low disk, usage hold)', the hold kinds (InboxDiskLow, InboxUsageReserve, InboxLoadHigh, InboxMemoryLow) come from the daemon, and internal/api/web_inbox.go's webInboxKinds deliberately leaves 'disk and usage holds' out of a project's web inbox. The plan does not say which holds count or how a machine-wide hold is attributed to one project (to every project, or none). The critic confirmed it: internal/api/inbox.go:760-763 inboxMachineWide treats every hold kind (disk-low, usage-reserve, load-high, memory-low) as the daemon's, the same in every project and listed once with no project. InboxItem.Project at line 130 is documented as absent for those items. internal/api/web_inbox.go:7-13 webInboxKinds leaves disk and usage holds out of a project's web view. The leaf counts 'holds' read 'from that project alone', and its fixture puts one hold in a project, but it never says how a machine-wide hold maps to a project. The count, and the guarantee that every surface shows the same number, are therefore undefined. Major holds for a contradiction. How should it be settled?

## Links

### Out

- blocks [web-project-needs-you](../capability/rhiza/web-monitor/web-project-needs-you.md)

### In

- answers [needs-you-count-excludes-machine-holds](../decision/needs-you-count-excludes-machine-holds.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
