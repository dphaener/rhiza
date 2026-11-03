# milestone-25

intent · `01M6HKVRE1B7RPP12EN6J2GGXC`

## milestone

25

## statement

The backlog is empty: every change set planned from the owner's queued requests and the backlog is built and shipped by the daemon on its own, so the owner can start on new things, bugs and polish

## Links

### Out

- decided-by [m25-empty-backlog](../decision/m25-empty-backlog.md)

### In

- serves change set `01M6HKTEKH6WQTCXRYDM9GSASG` (Milestone 25: empty the backlog — every planned change set from tonight's requests serves it so the daemon builds them)
- serves change set `01M6HHMTEDFK2AWQKX50EB77DD` (Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate)
- serves change set `01M6H3A65RDP761HKE82STZSQZ` (log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan))

## History

- `01M6HKTEKH6WQTCXRYDM9GSASG` Milestone 25: empty the backlog — every planned change set from tonight's requests serves it so the daemon builds them
