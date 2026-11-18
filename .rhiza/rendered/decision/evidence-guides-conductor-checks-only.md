# evidence-guides-conductor-checks-only

decision · `01M7PV72H3AAJ1X3210EC5TKCS`

## reasoning

internal/guide/text/conduct.md:38 has the conductor spawn only --role check-plan, check-code and check-qa. check-compliance is spawned only by classic kind=checks steps (internal/workflow/builtin/review.toml:45, .rhiza/workflows/docs-review.toml:31), which shallow-round-needs-leaf-coverage leaves unchanged with evidence optional. Telling the compliance check a round is refused in a conductor run describes a refusal it never meets, so the guide would be wrong; keeping it unchanged also matches the leaf's existing treatment of the other classic checks (check-quality, check-model in must_not). The alternative (keep it and explain a conductor might run it) has no support in the code or conduct guide, and evidence remains accepted-but-optional in classic rounds, so nothing is lost by leaving its guide as is.

## statement

Drop check-compliance from leaf check-guides-evidence-lines: the evidence-list and uncovered-leaf-refusal guidance goes only into the guides of the checks a conductor spawns (check-plan, check-code, check-qa) and their tests. internal/guide/text/check-compliance.md and internal/guide/check_compliance_test.go are left unchanged and move to the leaf's must_not, beside check-quality and check-model. The phrase 'check-compliance guidance used by conductors' in Decision shallow-round-needs-leaf-coverage is read as the three conductor check guides only; the refusal itself stays name-independent in code as that Decision says.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7pt83txevmbecpw6e6ee114-r3-f1](../question/review-01m7pt83txevmbecpw6e6ee114-r3-f1.md)

### In

- decided-by [check-guides-evidence-lines](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-guides-evidence-lines.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
