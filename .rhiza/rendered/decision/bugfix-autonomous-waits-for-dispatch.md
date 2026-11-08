# bugfix-autonomous-waits-for-dispatch

decision · `01M6WBX7DB8E3B3EY78RTNRTT0`

## reasoning

All three attempts failed for one reason, and retrying will not change it. The leaf's statement says the workflow is 'picked for bug requests by workflow dispatch', and its acceptance needs bugfix.toml mapped to kind bug. Dispatch is not on this branch or on origin/main (b696b22c has no internal/dispatch, and planRunArgs in internal/api/request.go:393 still picks only plan or plan-milestone). Attempt 2 already recorded this as question bugfix-needs-workflow-dispatch, and attempt 3 confirmed it. Builders that cannot write the code produce commits with only model records, which the verifier rightly rejects. The change set's intent names this end-to-end bug flow as its acceptance, so weakening or dropping the leaf would remove the change set's proof that it works. Merging the dispatch branch in would ship another change set inside this one, which the builder's note warned against. So the right answer is to sequence the leaf after dispatch, not to rebuild it now or reshape it.

## statement

Leaf bugfix-workflow-autonomous keeps its statement and acceptance unchanged and is not rebuilt until workflow dispatch (change set 01M6VY7G6HPJC3V0Q18E8EF9ZA) has landed on main. Then merge origin/main into cs/autonomy-per-workflow and build the leaf as specified: the testscript with a repo .rhiza/workflows/bugfix.toml (autonomous = true) mapped to kind bug, plan and signoff settings owner, the bug request going from plan run through gate and ship run with no owner inbox item. Do not merge cs/workflow-dispatch into this branch, and do not weaken the acceptance to pick the workflow some other way. Until dispatch lands the leaf stays blocked by bugfix-needs-workflow-dispatch; if this change set has to ship first, move the leaf into a follow-up change set that depends on dispatch rather than dropping it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6w86n-1](../question/leaf-failed-01m6w86n-1.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
