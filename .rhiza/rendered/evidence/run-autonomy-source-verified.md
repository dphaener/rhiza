# run-autonomy-source-verified

evidence · `01M6WBG2Y3FG260W7XKXBW544E` · actor [claude-code-19a95588](../actor/claude-code-19a95588.md)

## result

PASS

## statement

Verifier claude-code-19a95588 at b187f837: go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-autonomy-source|describe-golden)$' -v passes both. run-autonomy-source.txtar starts a run whose workflow declares [autonomy] plan = "auto" with rhiza.autonomy.accept = owner and accept undeclared; run start, run status show '- Autonomy plan: auto from workflow declares' and '- Autonomy accept: owner from rhiza.autonomy.accept'; run status --json carries stage/mode/source/workflow; the journal's started entry carries both. describe-golden.txtar is regenerated with the new autonomy field (RunStageAutonomy) and passes. The testscript's workflow is named declares, not plan, so I also ran a real plan run by hand in a scratch root (built binary, RHIZA_ env cleared, committed .rhiza/workflows/plan.toml = builtin plan + [autonomy] plan = "auto", git config rhiza.autonomy.accept owner): rhiza run start --workflow plan --new ... --dry-run printed 'Autonomy plan: auto from workflow plan' and 'Autonomy accept: owner from rhiza.autonomy.accept'. Diff read: runAutonomy takes the workflow declaration first, else rhiza.autonomy.<stage>; Started.Autonomy is folded into State and RunInfo; runLines prints one line per stage. It records all four stages (plan, accept, ship, signoff) rather than only those the run performs, which the acceptance allows. go test ./internal/run passes, and the internal/api tests for autonomy, run start and run status pass. The full internal/api package timed out under load average ~14 on 4 cores, with TestReadBudget failing on model read timings (show/ls/ready), which this change does not touch.

## Links

### Out

- verifies [run-autonomy-source](../capability/rhiza/workflows/workflow-runs/run-autonomy-source.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
