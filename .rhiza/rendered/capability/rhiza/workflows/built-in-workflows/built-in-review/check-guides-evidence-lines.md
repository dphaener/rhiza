# check-guides-evidence-lines

capability · `01M7PTHCYQHPSDSYPH4FEZ98AS` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/guide/... passes with new assertions in check_plan_test.go, check_code_test.go and check_qa_test.go that each guide's run findings --check example carries an "evidence" list with leaf and line, says one line per acceptance point, and says an uncovered leaf is refused; check_qa_test.go still asserts the Evidence node apply example.

## failure_waiver

Guide text only: no input path to refuse; the refusal itself is tested in check-round-leaf-coverage.

## files

internal/guide/text/check-plan.md, internal/guide/text/check-code.md, internal/guide/text/check-qa.md, internal/guide/check_plan_test.go, internal/guide/check_code_test.go, internal/guide/check_qa_test.go

## must_not

Change check-quality, check-model or check-compliance (internal/guide/text/check-compliance.md, internal/guide/check_compliance_test.go): classic checks steps keep evidence optional; change Go code outside internal/guide tests; drop the QA Evidence node instruction.

## statement

The check-plan, check-code and check-qa guides (the checks a conductor spawns) show the findings file with an evidence list and tell the check to write one evidence line per acceptance point of each leaf (what it ran or read, with the result), that a round leaving a leaf uncovered is refused in a conductor run, and check-qa still has the QA check write Evidence nodes on each leaf.

## Links

### Out

- decided-by [shallow-round-needs-leaf-coverage](../../../../../decision/shallow-round-needs-leaf-coverage.md)
- decided-by [evidence-guides-conductor-checks-only](../../../../../decision/evidence-guides-conductor-checks-only.md)

### In

- blocks [review-01m7pt83txevmbecpw6e6ee114-r3-f1](../../../../../question/review-01m7pt83txevmbecpw6e6ee114-r3-f1.md)
- verifies [check-guides-evidence-lines-verified](../../../../../evidence/check-guides-evidence-lines-verified.md)
- verifies [check-guides-evidence-lines-qa-1](../../../../../evidence/check-guides-evidence-lines-qa-1.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
