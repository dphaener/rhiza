# leaf-failed-01m6w86n-1

question · `01M6WBV4C8V31ZWE6WP87S9K1B`

## statement

Leaf bugfix-workflow-autonomous failed all 3 attempts of round 1: attempt 1 builder failed: no new commit since 54220972edbb; attempt 2 verifier failed: verifier found 1 problem(s) (f1 Acceptance not met at 2c1be359: the commit adds only model records (the blocking question bugfix-needs-workflow-dispatch and actor records), no code and no tes…); attempt 3 verifier failed: verifier found 1 problem(s) (f1 Acceptance not met at 5242ecef: the attempt-3 commit adds only model records (a note on question bugfix-needs-workflow-dispatch plus actor records), no code or…). How should it be built, or should it change?

## Links

### Out

- blocks [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

### In

- answers [bugfix-autonomous-waits-for-dispatch](../decision/bugfix-autonomous-waits-for-dispatch.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
