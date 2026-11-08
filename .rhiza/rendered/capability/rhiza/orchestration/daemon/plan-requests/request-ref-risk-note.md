# request-ref-risk-note

capability · `01M6W4ZN5TXW2QGS4G4047F38A` · parent [plan-requests](../plan-requests.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/request-ref-risk-note passes: with rhiza.autonomy.plan auto and rhiza.daemon.permissionMode bypassPermissions, request --ref queues the request (request ls shows it) and prints exactly one note line naming plan; the same request without --ref, and with rhiza.autonomy.plan owner, prints no note; describe-golden is regenerated

## statement

rhiza request --ref queues the request as before and, when the request-exposed list is not empty, its response carries a one-line note naming the flagged workflows and saying outside text reaches agents with no permission checks, printed as one line in human output; a request without --ref gets no note

## Links

### Out

- depends-on [request-exposed-workflows](../../../workflows/autonomous-runs/request-exposed-workflows.md)
- decided-by [risk-warnings-shape](../../../../../decision/risk-warnings-shape.md)

### In

- verifies [request-ref-risk-note-verified](../../../../../evidence/request-ref-risk-note-verified.md)
- verifies [request-ref-risk-note-qa-1](../../../../../evidence/request-ref-risk-note-qa-1.md)
- verifies [pr-365-risk-warnings](../../../../../evidence/pr-365-risk-warnings.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
