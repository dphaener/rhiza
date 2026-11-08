# stage-autonomy-from-workflow-verified-59ffe16b

evidence · `01M6WZCZV0GKFTXBDCHJXVMEAQ` · actor [claude-code-12b0dfdb](../actor/claude-code-12b0dfdb.md)

## result

PASS

## statement

Verified attempt 2 at 59ffe16b, the fix for the confirmed quality f2 / model f3 finding. Read the diff: daemonDemand in internal/api/daemon_projects.go now filters each ready candidate by workflowStageMode(ctx, ship.For(cs), "ship"), using the same daemonShipOf dispatch as daemonPoll1 (daemon.go:849-870), in place of the built-in ship-change declaration. Ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(daemon-root-ship-workflow|stage-autonomy|workflow-stage-autonomy|workflow-plan-autonomy|plan-workflow-autonomous|plan-workflow-autonomous-left|autonomy-accept-rule|bugfix-workflow-autonomous|run-start-ready-dispatch|run-autonomy-source|daemon-root.*)$': ok. These cover the acceptance (plan autonomous/owner overrides, accept = owner leaves drafts proposed, ship = auto starts while the setting is owner, signoff = owner starts without --autonomous, a ship workflow's plan ignored, stage-autonomy unchanged). Ran the new daemon-root-ship-workflow.txtar against 59ffe16b~ in a scratch worktree: it fails ('no match for ^Started a run of'), so it covers the finding. The full internal/api package did not finish here: one run hit the 600s timeout at load ~15, one tripped the known write into ~/code's root state folder, and a scratch-worktree rerun was stopped for time. CI runs it before merge.

## Links

### Out

- verifies [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
