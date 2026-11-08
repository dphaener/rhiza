# risk-warnings-shape

decision · `01M6W4X6742B6C7HXEKS1W4PPT`

## reasoning

Read the question, the change set's intent (warn, never block; follows workflows-are-the-rules and warn-dont-safeguard from PR #353), the Workflow dispatch change set (01M6VY7G6HPJ, accepted, unbuilt; decision workflow-dispatch-shape serves the backlog milestone 01M6HKVR, milestone-25), and the landed Scheduling as policy (#358), which ships backlog work by default, so a milestone is not needed to get built. Its sibling's milestone is the obvious home; milestone-29 is about one root running several projects, which this does not touch. The code: request.go stores --ref as opaque, untrusted text and passes it into the plan run; autonomy is repo-wide git config (daemon_settings.go rhiza.autonomy.*); the daemon starts runs with rhiza.daemon.permissionMode for steps that set none (daemon.go), and workflow steps may set permission_mode (workflow.go); adapter.Bypass is the existing test for a mode that skips permission checks, and bypassCheck already gates it on a machine with trust = dedicated, so a warning is the only thing left to add. Because Rhiza cannot tell a cron-queued request from a hand-typed one, any request has to count as outside input, and gating on the autonomy stages keeps the warning off repos where a person reads the plan before anything runs. health.go already has an autonomy part listing those settings, and warn is the state for 'worth a look' that never makes health exit 1, which fits never-blocks. Waiting on dispatch would hold this back for a dependency it only extends; a single reaches-from-request function lets either order work.

## statement

(1) Risk warnings serves @milestone-25 (the backlog milestone), like its sibling Workflow dispatch: it is polish on request intake, not multi-project work. (2) Outside input means a queued request: its text and --ref are outside text, and a schedule is just a cron job running rhiza request, which Rhiza cannot see, so every request counts. A workflow is flagged when the daemon would run it from a request with no person in between and any agent session it starts, following its included workflows (agent and arbiter steps, a leaves step's builder and verifier, and a checks step's checks, critic, fixer and verifier), runs in a bypass mode. Concretely: the request workflow (plan today) when rhiza.autonomy.plan is on, and the ready workflow (ship-change today) when rhiza.autonomy.plan, rhiza.autonomy.accept and rhiza.autonomy.ship are all on, because only then does request text reach its agents unseen (with accept owner a person accepts the planned change set before the ready workflow can run). A session's effective mode is its step's or role's own permission_mode, else the run's (the repo's rhiza.daemon.permissionMode), and bypass is whatever the agent's adapter says it is (adapter Bypass: bypassPermissions for Claude Code). rhiza workflow check <name> reports this as a warning, beside its problems: valid stays true and the exit code does not change. (3) In rhiza health the warning goes in the autonomy part: state warn (never bad, so health never exits 1 over it), the summary gains 'outside text reaches agents with no permission checks', there is one detail line per flagged workflow naming it, its trigger and the bypassing step(s), and Next is rhiza guide workflows. (4) Do not wait for Workflow dispatch. Work out which workflows a request reaches in one function that today returns the fixed picks (request: plan; ready: ship-change, or the change set's own workflow). Whichever of the two change sets lands second has that function read .rhiza/dispatch.toml (every [request] and [ready] rule plus the built-in defaults), so rules added by dispatch get flagged too.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-0381a0e6-1](../question/question-0381a0e6-1.md)
- derived-from [risk-warnings-plan](../source/risk-warnings-plan.md)

### In

- decided-by [request-exposed-workflows](../capability/rhiza/workflows/autonomous-runs/request-exposed-workflows.md)
- decided-by [workflow-check-risk-warning](../capability/rhiza/workflows/workflow-authoring/workflow-check-risk-warning.md)
- decided-by [health-autonomy-risk](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-autonomy-risk.md)
- decided-by [request-ref-risk-note](../capability/rhiza/orchestration/daemon/plan-requests/request-ref-risk-note.md)
- decided-by [workflows-guide-outside-input](../capability/rhiza/workflows/workflow-authoring/workflows-guide-outside-input.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
