# check-round-leaf-coverage

capability · `01M7PTHCNE0M4REEQY57T5PW8Z` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/api/... ./internal/run/... passes with new tests in internal/api/run_review_test.go: in a conductor step a round with one evidence line per leaf is journaled with its lines on the round; a round whose only finding names leaf A and evidence covers leaf B is accepted; a change set with no leaves accepts {"findings":[]}; at a classic checks step {"findings":[]} is still accepted; (failure case) in a conductor step {"findings":[]} with two leaves is a usage error naming both leaves with Next pointing at the run.findings schema and the journal holds no new entry; (failure case) an evidence line whose leaf is not a leaf of the change set, or whose line is blank, is a usage error. go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after regenerating it with -update, and make self-check passes (docs/reference/cli.md regenerated).

## files

internal/api/run_review.go, internal/api/run_review_test.go, internal/api/run_round_recorded_test.go, internal/api/session_done_check_test.go, internal/api/skill_check_test.go, internal/run/journal.go, internal/run/journal_test.go, cmd/rhiza/testdata/script/describe-golden.txtar, docs/reference/cli.md

## must_not

Key any rule on a check's name (plan, code, qa); add a time or duration heuristic; change plan-review rounds, --leaf verifier findings or classic checks steps beyond accepting an optional evidence list; change rhiza run review-check; edit cmd/rhiza testscripts other than describe-golden.txtar.

## statement

In a conductor step (workflow.Step.Conducts()), rhiza run findings --check takes an optional evidence list ({leaf, line}) beside findings and refuses, journaling nothing, a round that leaves any live leaf of the run's change set uncovered (no finding naming it and no non-blank evidence line for it), naming each uncovered leaf; evidence lines are journaled on the findings entry and round.

## Links

### Out

- decided-by [shallow-round-needs-leaf-coverage](../../../../../decision/shallow-round-needs-leaf-coverage.md)

### In

- depends-on [check-round-evidence-status](check-round-evidence-status.md)
- depends-on [conduct-fixtures-cover-leaves](conduct-fixtures-cover-leaves.md)
- verifies [check-round-leaf-coverage-verified](../../../../../evidence/check-round-leaf-coverage-verified.md)
- verifies [check-round-leaf-coverage-qa-1](../../../../../evidence/check-round-leaf-coverage-qa-1.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
