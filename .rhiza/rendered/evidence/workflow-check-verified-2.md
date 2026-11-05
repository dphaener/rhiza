# workflow-check-verified-2

evidence · `01M6KQ0DDY9FY6N5YMRZYPQK5C` · actor [claude-code-ea125804](../actor/claude-code-ea125804.md)

## result

PASS

## statement

Verifier, attempt 2 at 268052e3. go test -count=1 ./cmd/rhiza -run TestScripts/workflow-check passes. The script runs workflow check on the built-ins build, plan, review and ship-change and on copies of internal/workflow/testdata (same bytes as valid, bad-outcome, missing-target and unknown-kind): valid.toml and every built-in exit 0. bad-outcome, missing-target and unknown-kind each exit non-zero naming the step and the problem. many.toml reports all 3 problems. A missing include is reported with the step that includes it and is expanded once the repo layer holds it. A name that is nowhere is a usage error. --json is cmp-equal to rhiza api call workflow.check for a bad and a good workflow. Read the fix diff (agentProblems in internal/api/workflow_check.go): it runs run start's adapter lookup and startAdapter mode checks on each agent and permission_mode a step sets, for agent, arbiter, leaves and checks roles. Reproduced the earlier quality finding f3 with a scratch build in a scratch repo: a repo workflow whose agent step sets agent = "nope" now exits 1 with 'step build: the agent step: unknown agent "nope"', agent = "claude-code" stays valid with exit 0, and --json matches api call byte for byte. go test ./internal/workflow ok. go test -timeout 25m ./internal/api ok (673s; a first run hit the default 10m timeout on a busy machine)

## Links

### Out

- verifies [workflow-check](../capability/rhiza/workflows/workflow-authoring/workflow-check.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
