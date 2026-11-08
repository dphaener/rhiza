# request-exposed-workflows-verified

evidence · `01M6W73T7D4BDFR78YNF7AVHB4` · actor [claude-code-6a2e52c0](../actor/claude-code-6a2e52c0.md)

## result

PASS

## statement

At 0e0a5b80 I ran go test -trimpath -count=1 ./internal/api -run TestRequestExposedWorkflows -v: all 10 subtests PASS; go vet and gofmt are clean on internal/api. I read internal/api/request_exposure.go and its test. requestExposedWorkflows picks plan (trigger request) always and ship-change (trigger ready) only when rhiza.autonomy.plan, accept and ship are all auto. It resolves each through resolveWorkflow, which inlines included workflows, so steps come out as e.g. build/build builder. bypassSteps checks agent steps, leaves builder/verifier, checks-step checks/critic/fixer/verifier and the arbiter, each with its own permission_mode, else rhiza.daemon.permissionMode, against its adapter's Bypass. The tests cover every acceptance point: plan auto under bypass flags plan only; plan+accept+ship auto flags plan and ship-change; accept owner flags plan only; step acceptEdits under a bypass daemon mode is not flagged, and a step bypassPermissions with no daemon mode is; the arbiter is counted; ship-change is flagged through an included builder alone and an included critic alone, naming that step and role; entries carry workflow, trigger and steps. The acceptance's 'plan owner flags nothing' case is tested as 'plan owner still flags plan, not ship-change', per decision request-workflow-exposure-ignores-plan-stage, which answers blocking question review-01m6w4k1-r3-f1 and rewrites that case. There is no .rhiza/dispatch.toml or dispatch support in this tree, so the dispatch clause does not apply yet.

## Links

### Out

- verifies [request-exposed-workflows](../capability/rhiza/workflows/autonomous-runs/request-exposed-workflows.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
