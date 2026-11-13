# needs-you-count-excludes-machine-holds

decision · `01M7B2MB5BAH9EW1YHAKWVJNWV`

## reasoning

Every hold kind in the code is the daemon's, not a project's: internal/api/inbox.go inboxMachineWide treats disk-low, usage-reserve, load-high and memory-low as machine-wide, listed once with no project, and InboxItem.Project is documented as absent for them. No per-project hold kind exists. The web UI already settled this for a project's view: internal/api/web_inbox.go webInboxKinds leaves disk and usage holds out of a project's web inbox. Attributing a machine-wide hold to every project would inflate each count with the same item and break 'read from that project alone'; attributing it to one is arbitrary. Excluding holds keeps the count per-project and the same on every surface, and matches the owner-facing web inbox. Showing machine-wide holds elsewhere in the frame is outside this question and this change set's intent (a project picker with each project's needs-me count).

## statement

A project's needs-you count holds no holds. The daemon's machine-wide holds (disk-low, usage-reserve, load-high, memory-low) belong to no project, so they are counted in no project's needs-you count and on no surface that shows that count (header, picker, Home's title). The count is open questions, change sets waiting for the owner's sign-off, failed runs and stuck runs, read from that project alone. web-project-needs-you's fixture keeps its one hold to prove it adds nothing: the count for that fixture is 4, not 5.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f2](../question/review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f2.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
