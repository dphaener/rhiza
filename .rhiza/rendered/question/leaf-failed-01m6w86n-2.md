# leaf-failed-01m6w86n-2

question · `01M6WCFVHRN3QSVM50RKQ9A8VP`

## statement

Leaf bugfix-workflow-autonomous failed all 3 attempts of round 2: attempt 4 builder failed: no new commit since 5242ecefc99a; attempt 5 verifier failed: verifier found 1 problem(s) (f1 Acceptance not met at 36d24c64: the attempt-5 commit only adds model notes (.rhiza/log JSON) recording that workflow dispatch (change set 01M6VY7G6HPJC3V0Q18E8…); attempt 6 verifier failed: verifier found 1 problem(s) (f1 Acceptance not met at 3e4c472d: the attempt-6 commit adds only .rhiza/log model records (a note that workflow dispatch is still unlanded plus actor records), n…). How should it be built, or should it change?

## Links

### Out

- blocks [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
