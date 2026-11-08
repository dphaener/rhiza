# request-ref-risk-note-verified

evidence · `01M6W7RRJDX12H9MV0E1MJJQJ6` · actor [claude-code-9401a472](../actor/claude-code-9401a472.md)

## result

PASS

## statement

At 601a755d, go test -trimpath ./cmd/rhiza -run 'TestScripts/(request-ref-risk-note|describe-golden)$' -count=1 -v: both PASS. The script sets rhiza.autonomy.plan auto and rhiza.daemon.permissionMode bypassPermissions, runs request --ref, sees 'Queued feature request', exactly one '^Note: .*workflow plan' line saying 'no permission checks', request ls lists the request, and --json carries a note naming plan; a request without --ref and one with plan owner print no Note. describe-golden passes with the new note field. Read the diff: runRequest still queues as before and only adds requestRiskNote, which returns empty for no ref, unreadable settings or an empty requestExposedWorkflows list.

## Links

### Out

- verifies [request-ref-risk-note](../capability/rhiza/orchestration/daemon/plan-requests/request-ref-risk-note.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
