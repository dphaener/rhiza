# review-01m7pt83txevmbecpw6e6ee114-r3-f1

question · `01M7PV67Z9271HEMW6WVBX719P`

## statement

Plan review round 3 left finding f1 (contradiction on check-guides-evidence-lines): The leaf (and Decision shallow-round-needs-leaf-coverage) treat check-compliance as guidance used by conductors, but internal/guide/text/conduct.md:38 spawns only --role check-plan, check-code and check-qa; check-compliance is the classic review workflow's check (kind checks), which the Decision leaves unchanged. Telling the compliance check that an uncovered leaf is refused in a conductor run describes a refusal it never meets; drop check-compliance.md and check_compliance_test.go from the leaf, or say why a conductor would run it. The critic confirmed it: Real: internal/guide/text/conduct.md:38 has the conductor spawn only --role check-plan, check-code and check-qa, and review.md:7-13 lists check-compliance as the classic workflows' check (docs-review.toml declares it under a kind = "checks" step). The Decision keeps classic checks steps unchanged with evidence optional, so a compliance check never meets the uncovered-leaf refusal; the leaf telling it that such a round is refused in a conductor run, and the Decision's 'check-compliance guidance used by conductors', contradict the code. Major for a contradiction stands. How should it be settled?

## Links

### Out

- blocks [check-guides-evidence-lines](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-guides-evidence-lines.md)

### In

- answers [evidence-guides-conductor-checks-only](../decision/evidence-guides-conductor-checks-only.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
