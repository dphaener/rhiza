# plan-clean-round-skips-critic-step

capability · `01M77NZCXG646BQB0FY09EKQPW` · parent [plan-workflow](../plan-workflow.md)

## acceptance

New tests in internal/api and/or internal/workflow, plus a testscript under cmd/rhiza/testdata/script: in a plan run whose review visit journals {"findings":[]}, the run moves from review to fresh and check (and on to gate) without starting a critic session; a round with findings still goes to the critic; both built-in plan and plan-milestone workflows route this way (a workflow test reads both files); rhiza run verdict with {"verdicts":[]} still exits non-zero. Run: go test -trimpath ./internal/api/... ./internal/workflow/... -run <the new tests>, and go test -trimpath ./cmd/rhiza -run TestScripts/<the new script>; regenerate describe-golden if a request type changes.

## statement

In the plan workflows (plan and plan-milestone), after the review step the run goes to the critic only when the latest review round holds at least one finding; a clean round (no findings) goes straight to fresh, then check, with no critic session, and rhiza run verdict still refuses an empty verdicts list.

## Links

### Out

- decided-by [plan-clean-round-skips-critic](../../../../../decision/plan-clean-round-skips-critic.md)
- derived-from [plan-clean-round-critic-plan](../../../../../source/plan-clean-round-critic-plan.md)

### In

- verifies [plan-clean-round-skips-critic-step-verified](../../../../../evidence/plan-clean-round-skips-critic-step-verified.md)
- verifies [plan-clean-round-skips-critic-step-qa-1](../../../../../evidence/plan-clean-round-skips-critic-step-qa-1.md)
- verifies [plan-clean-round-skips-critic-step-verified-2](../../../../../evidence/plan-clean-round-skips-critic-step-verified-2.md)
- verifies [pr-486-fix-plan-workflows-visit](../../../../../evidence/pr-486-fix-plan-workflows-visit.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
