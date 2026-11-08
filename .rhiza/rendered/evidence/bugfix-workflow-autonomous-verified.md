# bugfix-workflow-autonomous-verified

evidence · `01M6WQYT62S0FH0GAXYM5RMY0D` · actor [claude-code-230555bd](../actor/claude-code-230555bd.md)

## result

PASS

## statement

At 7a89e838 ran `go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/bugfix-workflow-autonomous$'`: ok (124s). The script uses fake rex, claude and gh; the repo has .rhiza/workflows/bugfix.toml and bugfix-ship.toml with autonomous = true, dispatch.toml maps request.bug=bugfix and ready.bug=bugfix-ship, and git config rhiza.autonomy.plan and rhiza.autonomy.signoff are owner. It asserts the daemon logs the bug request's plan run start with --autonomous (planRunAutonomous reads the workflow's declaration before the setting), the plan run succeeds with Autonomous: yes, gate pass and one accepted event, the change set is accepted by rule autonomous, the daemon logs the ship run start with --workflow bugfix-ship --pick rule ready.bug --autonomous (daemonStart drops --autonomous when the signoff stage is owner, so the declaration is what keeps it), the ship run succeeds signed off by rule autonomous with no actor, and rhiza inbox --json holds no item after each phase. Read the diff: test-only plus the builder's actor record.

## Links

### Out

- verifies [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
