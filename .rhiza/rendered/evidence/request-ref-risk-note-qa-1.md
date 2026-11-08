# request-ref-risk-note-qa-1

evidence · `01M6W9QE38RXP1QT1ZEMXZAKHR` · actor [claude-code-afe51521](../actor/claude-code-afe51521.md)

## result

PASS

## statement

QA at bbfe87e3 with a built binary (go build -trimpath) in a scratch Rhiza root and project, real git config, no fakes: With plan auto + bypassPermissions: rhiza request --ref https://example.com/i/1 'Add a thing' exits 0, prints 'Queued feature request ...' then one 'Note: outside text reaches agents with no permission checks in workflow plan, ship-change (rhiza guide workflows)'; --json has the note field; request ls lists the queued requests; a request without --ref prints no note; --dry-run with --ref notes and queues nothing. plan owner + acceptEdits with --ref: no note.

## Links

### Out

- verifies [request-ref-risk-note](../capability/rhiza/orchestration/daemon/plan-requests/request-ref-risk-note.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
