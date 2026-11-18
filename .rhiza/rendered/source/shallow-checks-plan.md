# shallow-checks-plan

source · `01M7PTGHZ2BJK9MG694YJ9S3GK`

## quotes

```json
[]
```

## raw_sha256

667f398e30778327ab60df623a38be0a984bcc9ac9458e140d8f90047c555336

## session

e50b77ed-cc5e-4e7f-950d-01a1eda447fd

## summary

Autonomous plan run for the bug that a conductor run's check can record an empty findings round with no evidence and pass. Research: runCheckFindings in internal/api/run_review.go accepts any empty list; checks are identified only by free-form names (plan, code, qa in ship-change); workflow.Step.Conducts() marks an agent step declaring checks. The planner asked one design question (how to refuse a shallow round); there was no owner, and the arbiter ruled in Decision shallow-round-needs-leaf-coverage: in a conductor step every leaf must be covered by a finding naming it or a non-blank evidence line, else run findings --check refuses and journals nothing; evidence is journaled and shown by run status; classic checks steps and plan review unchanged; the check guides write one evidence line per acceptance point. The plan splits this into the refusal in run findings, run status showing the lines, the conductor testscript fixtures, and the guides. The change set serves m36-conductor.

## title

Plan of Fix Shallow checks pass (autonomous plan run)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
