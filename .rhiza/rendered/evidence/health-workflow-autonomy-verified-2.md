# health-workflow-autonomy-verified-2

evidence · `01M6WQY2B0VNP49KJMRGY7DR81` · actor [claude-code-d57d46ca](../actor/claude-code-d57d46ca.md)

## result

PASS

## statement

Verifier claude-code-d57d46ca at ab0aeb30: ran go test -trimpath ./cmd/rhiza -run 'TestScripts/^health$' -count=1, PASS. health.txtar checks that with .rhiza/workflows/bugfix.toml declaring autonomous = true the autonomy detail lists 'workflow bugfix (.rhiza/workflows/bugfix.toml): plan auto, accept auto, ship auto, signoff auto' and 'workflow plan-milestone (builtin:plan-milestone): plan owner, accept owner, ship owner, signoff owner'; with no repo declaration only plan-milestone is listed and no bugfix; health --json and api call health compared with cmp in both states. The leaf commit adds only model records (confirmation that the capability was already built).

## Links

### Out

- verifies [health-workflow-autonomy](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-workflow-autonomy.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
