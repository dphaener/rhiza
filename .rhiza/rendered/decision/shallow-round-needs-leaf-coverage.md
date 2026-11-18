# shallow-round-needs-leaf-coverage

decision · `01M7PTF00HA0WBMQ27P1K6W4BT`

## reasoning

Read the question, the change set intent, internal/api/run_review.go runCheckFindings (accepts any empty list today, which is the bug), internal/workflow/workflow.go (Check is only Name+Role; Conducts() = agent step with checks) and builtin ship-change.toml (plan/code/qa are bare names). Because check names are free config, the refusal has to be name-independent; per-leaf coverage is the common denominator of the intent's two examples (QA: Evidence per leaf; plan: a line per acceptance point) and it maps onto data the command already resolves (changeLeafResolver). Counting 'acceptance points' inside free-text acceptance fields would be fragile, so the gate requires at least one line per leaf and the skills ask for one per point. Letting a finding on a leaf count as coverage avoids forcing redundant lines, and requiring coverage of every leaf (not only on empty rounds) stops a shallow check from slipping past with a single nit. Requiring Evidence nodes for every check would wrongly burden plan/code checks, and a minimum-time gate is model- and machine-dependent and flaky. Refusing at record time (rather than in review-check) makes the conductor rerun the check immediately, as the intent asks, and keeps review-check's semantics. Scoping to Conducts() steps honours the question's 'classic checks steps unchanged'.

## statement

In a conductor step (workflow.Step.Conducts(): an agent step declaring checks), `rhiza run findings --check` refuses a round unless every leaf of the run's change set is covered, whatever the check is named: the findings file gains an `evidence` list, `{"findings":[...],"evidence":[{"leaf":<handle or ID>,"line":<what the check ran or read for that leaf's acceptance, with its result>}]}`, and a leaf counts as covered when the round has a finding naming it or at least one non-blank evidence line for it. An uncovered leaf (including a clean `{"findings":[]}` with no evidence) is a usage error naming each uncovered leaf, with Next pointing at the schema; nothing is journaled, so the check has recorded no round and the conductor reruns it. Evidence lines are journaled on the findings entry and shown by run status, so the conductor and critic can read what the check did. An evidence leaf must resolve to one of the change set's leaves (usage error otherwise). A change set with no leaves needs no evidence. No per-name rules (no special case for plan/code/qa) and no minimum-duration heuristic. Classic checks steps (kind checks) and plan-review rounds are unchanged: evidence is optional there and ignored when absent. The check skills (check-plan, check-code, check-qa, check-compliance guidance used by conductors) are updated to write one evidence line per acceptance point per leaf; a QA check also still writes Evidence nodes on leaves as today.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-e50b77ed-1](../question/question-e50b77ed-1.md)

### In

- decided-by [check-round-leaf-coverage](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-leaf-coverage.md)
- decided-by [check-round-evidence-status](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-evidence-status.md)
- decided-by [conduct-fixtures-cover-leaves](../capability/rhiza/workflows/built-in-workflows/built-in-review/conduct-fixtures-cover-leaves.md)
- decided-by [check-guides-evidence-lines](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-guides-evidence-lines.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
