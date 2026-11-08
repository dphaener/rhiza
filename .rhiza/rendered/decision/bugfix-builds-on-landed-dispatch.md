# bugfix-builds-on-landed-dispatch

decision · `01M6WKGN81K7SZ36RB0PJEADB1`

## reasoning

The question asked to land dispatch and then rebuild the leaf on it. Dispatch merged to main as PR #367 and origin/main was merged into this branch (merge commit ea0899ea), which resolved the overlap so a request's plan run takes the workflow dispatch picks and autonomy reads that workflow's declaration. Nothing still blocks the leaf, so the retry builds it as planned.

## statement

Build leaf bugfix-workflow-autonomous on workflow dispatch as it landed: dispatch (change set 01M6VY7G6HPJC3V0Q18E8EF9ZA, PR #367) is on main and has been merged into cs/autonomy-per-workflow, so a bugfix.toml mapped to request kind bug through .rhiza/dispatch.toml is reachable and the leaf's acceptance stands unchanged.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [bugfix-needs-workflow-dispatch](../question/bugfix-needs-workflow-dispatch.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
