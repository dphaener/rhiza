# pr-365-risk-warnings

evidence · `01M6WAYVH1BPZFWD6GQWFY2XME` · actor [dphaener](../actor/dphaener.md)

## head

9813a05ed808e1aff7f0e4c165455dd4ad7f5803

## pr

https://github.com/dphaener/rhiza/pull/365

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #365 (https://github.com/dphaener/rhiza/pull/365) passed CI at 9813a05

## status

verified

## Links

### Out

- verifies [request-exposed-workflows](../capability/rhiza/workflows/autonomous-runs/request-exposed-workflows.md)
- verifies [workflow-check-risk-warning](../capability/rhiza/workflows/workflow-authoring/workflow-check-risk-warning.md)
- verifies [health-autonomy-risk](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-autonomy-risk.md)
- verifies [request-ref-risk-note](../capability/rhiza/orchestration/daemon/plan-requests/request-ref-risk-note.md)
- verifies [workflows-guide-outside-input](../capability/rhiza/workflows/workflow-authoring/workflows-guide-outside-input.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
