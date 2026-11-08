# health-workflow-autonomy-verified

evidence · `01M6WDY9MNEQGGA49Q9RK4Z90E` · actor [claude-code-c778e88b](../actor/claude-code-c778e88b.md)

## result

PASS

## statement

At 1bf699e7: go test -trimpath ./cmd/rhiza -run 'TestScripts/^health$' passes; health.txtar asserts the built-in plan-milestone line alone, then with .rhiza/workflows/bugfix.toml (autonomous = true) the line 'workflow bugfix (.rhiza/workflows/bugfix.toml): plan auto, accept auto, ship auto, signoff auto' plus plan-milestone with plan, accept, ship and signoff owner, and cmp of health --json against api call health. By hand with a scratch build in a fresh repo: without a repo workflow the autonomy detail lists only 'workflow plan-milestone (builtin:plan-milestone): plan owner, accept owner, ship owner, signoff owner' after the four rhiza.autonomy.* settings; with bugfix.toml added it also lists bugfix with its file and all four stages auto; health --json and api call health are byte-identical. Diff reviewed: healthWorkflowAutonomy unions personal, repo and built-in names, resolves each through resolveWorkflow (so precedence holds), and lists only workflows declaring a stage. go test -short -run Health ./internal/api passes; the other internal/api failures seen (TestReadBudget timeout at load 15, a usage.live.json leak into the enclosing root, TestRunStartAgentRefusalsBeforeWrite) are unrelated to health, and the last also fails at the parent 1e97bc3f.

## Links

### Out

- verifies [health-workflow-autonomy](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-workflow-autonomy.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
