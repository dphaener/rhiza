# workflow-check-risk-warning-verified

evidence · `01M6W7BBAWMHXAAYEA16CW57JM` · actor [claude-code-b34a0c3a](../actor/claude-code-b34a0c3a.md)

## result

PASS

## statement

At c849938f I ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(workflow-check-risk|describe-golden|workflow-check)' and go test -trimpath -count=1 ./internal/api -run 'TestExposureWarnings|TestRequestExposedWorkflows|WorkflowCheck': both ok; gofmt and go vet clean. The workflow-check-risk script sets rhiza.autonomy.plan auto and rhiza.daemon.permissionMode bypassPermissions and checks workflow check plan stays valid with a human 'warning:' line and a JSON warnings field, no warning without bypass mode, and none for build. describe-golden carries the new warnings field. By hand, with a scratch build in a scratch repo: plan auto + bypass gives exit 0, valid true, and 'workflow plan runs from a queued request with bypassed permissions in interview, review, critic, revise, arbiter: outside text reaches agents with no permission checks' in both human and JSON output; plan, accept and ship auto gives ship-change the same warning naming 'a ready change set' and its bypassing steps, exit 0; with no permission mode there is no warning. With rhiza.autonomy.plan owner plus bypass the warning still prints. The acceptance's 'plan owner prints no warning' case is replaced by decision request-workflow-exposure-ignores-plan-stage (plan owner with a bypass daemon mode flags plan), which request-exposed-workflows-verified already applied, and the builder's commit notes this. I read the diff: warnings is omitempty, it is computed only for named (not path) workflows, it never touches valid or the exit code, and Human prints it after the verdict or problems.

## Links

### Out

- verifies [workflow-check-risk-warning](../capability/rhiza/workflows/workflow-authoring/workflow-check-risk-warning.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
