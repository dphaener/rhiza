# run-autonomy-source-verified-5eb92da1

evidence · `01M6WN6FY75E7FFYERPZXJHHHA` · actor [claude-code-7712db35](../actor/claude-code-7712db35.md)

## result

PASS

## statement

Verifier claude-code-7712db35 at 5eb92da1 (builder commit only tracks its actor record; the code landed in b187f837): go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-autonomy-source|describe-golden|workflow-plan-autonomy)$' passes. run-autonomy-source.txtar starts a run of a workflow declaring [autonomy] plan = "auto" with accept undeclared and git config rhiza.autonomy.accept owner; run start and run status print '- Autonomy plan: auto from workflow declares' and '- Autonomy accept: owner from rhiza.autonomy.accept'; run status --json carries stage/mode/source/workflow; the journal's started entry carries both. describe-golden.txtar includes the autonomy field and passes. Since the script's workflow is a command workflow named declares, I also ran a real plan run by hand with a scratch build in a scratch root (RHIZA_* env cleared): .rhiza/workflows/plan.toml = builtin plan.toml + [autonomy] plan = "auto", rhiza.autonomy.accept owner; rhiza run start --workflow plan --new Exports --dry-run --permission-mode manual printed 'Autonomy plan: auto from workflow plan' and 'Autonomy accept: owner from rhiza.autonomy.accept' (plus ship/signoff from settings). Read internal/api/autonomy.go runAutonomy: workflow declaration first, naming the workflow, else rhiza.autonomy.<stage>; recorded in run.Started.Autonomy.

## Links

### Out

- verifies [run-autonomy-source](../capability/rhiza/workflows/workflow-runs/run-autonomy-source.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
