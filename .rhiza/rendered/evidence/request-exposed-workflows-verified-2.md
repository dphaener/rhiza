# request-exposed-workflows-verified-2

evidence · `01M6WA2FCZ26JYKGMVKHRR74DP` · actor [claude-code-ed85aa75](../actor/claude-code-ed85aa75.md)

## result

PASS

## statement

Attempt 2 at d35520fc (a model-only commit rewriting the leaf's statement and acceptance to match decision request-workflow-exposure-ignores-plan-stage). Ran go test -trimpath ./internal/api -run TestRequestExposedWorkflows -v: PASS. Read the tests and internal/api/request_exposure.go against each acceptance point: plan auto under bypassPermissions flags plan only; plan/accept/ship auto flags plan (request) and ship-change (ready); accept owner flags plan only; plan owner with ship auto under bypass flags plan, not ship-change; a step acceptEdits under a bypass daemon mode is not flagged, a step bypassPermissions with no daemon mode is; an arbiter step with bypassPermissions is counted; ship-change is flagged through an included build builder alone ('build/build builder') and an included review critic alone ('review/review critic') with no daemon mode; each entry carries workflow, trigger and steps. No dispatch.toml support exists in this code, so the conditional dispatch clause does not apply. The rewritten statement now matches requestExposedWorkflows (plan always picked; ship-change only when plan, accept and ship are all auto). make self-check reports only stale-render failures (rhiza render runs after landing), no model problems.

## Links

### Out

- verifies [request-exposed-workflows](../capability/rhiza/workflows/autonomous-runs/request-exposed-workflows.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
