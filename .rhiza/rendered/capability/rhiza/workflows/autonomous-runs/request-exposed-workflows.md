# request-exposed-workflows

capability · `01M6W4ZMW57BSXGYJYXV6XPWP2` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

Unit tests in internal/api (go test -trimpath ./internal/api -run TestRequestExposedWorkflows) cover: plan auto with rhiza.daemon.permissionMode bypassPermissions flags plan only; plan, accept and ship auto flags plan and ship-change; plan and ship auto with rhiza.autonomy.accept owner flags plan only; plan owner with ship auto under a bypass daemon mode flags plan, not ship-change; a step permission_mode acceptEdits under a bypass daemon mode is not flagged and a step permission_mode bypassPermissions with no daemon mode is; arbiter steps are counted; ship-change is flagged through an included workflow's role alone (a leaves-step builder or verifier, or a checks-step check, critic or fixer, with permission_mode bypassPermissions and no daemon mode), naming that step and role; each entry names the workflow, its trigger (request or ready) and the bypassing steps; and, when dispatch.toml support exists, a workflow named by a [request] or [ready] rule is flagged

## statement

One function in internal/api lists the workflows the daemon would run from a queued request with no person in between, each with its trigger and the steps and roles whose effective mode the session's adapter reports as a bypass mode, counting every agent session the workflow starts, following its included workflows: agent and arbiter steps, a leaves step's builder and verifier, and a checks step's checks, critic, fixer and verifier, each with its own permission_mode, else the run's (the repo's rhiza.daemon.permissionMode): the request workflow (plan) whatever rhiza.autonomy.plan says, milestone requests included (an interactive plan run is no gate: its agents read the request text before the owner joins), and the ready workflow (ship-change) when rhiza.autonomy.plan, rhiza.autonomy.accept and rhiza.autonomy.ship are all auto (with accept owner a person accepts the planned change set before it is ready); where .rhiza/dispatch.toml exists in the code it builds on (Workflow dispatch landed first), it reads every [request] and [ready] rule plus the built-in defaults instead of the fixed picks

## Links

### Out

- decided-by [risk-warnings-shape](../../../../decision/risk-warnings-shape.md)

### In

- depends-on [workflow-check-risk-warning](../workflow-authoring/workflow-check-risk-warning.md)
- depends-on [health-autonomy-risk](../../orchestration/daemon/daemon-status/health-command/health-autonomy-risk.md)
- depends-on [request-ref-risk-note](../../orchestration/daemon/plan-requests/request-ref-risk-note.md)
- blocks [review-01m6w4k1-r3-f1](../../../../question/review-01m6w4k1-r3-f1.md)
- verifies [request-exposed-workflows-verified](../../../../evidence/request-exposed-workflows-verified.md)
- verifies [request-exposed-workflows-qa-1](../../../../evidence/request-exposed-workflows-qa-1.md)
- verifies [request-exposed-workflows-verified-2](../../../../evidence/request-exposed-workflows-verified-2.md)
- verifies [pr-365-risk-warnings](../../../../evidence/pr-365-risk-warnings.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
