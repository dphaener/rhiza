# plan-reviewer-leaf-contract-verified

evidence · `01M7KSNGK2VMPCYB48YZQRBTFJ` · actor [claude-code-0b1443d5](../actor/claude-code-0b1443d5.md)

## result

PASS

## statement

Verified at fd9a92409: go test -trimpath -p 2 ./internal/guide/... -count=1 ok; TestPlanReviewerText (eight kinds, no-failure-case after oversized-leaf, missing-acceptance narrowed to vague/unrunnable/full-suite, Must not names rhiza run review-check), TestPlanReviserText and TestPlanCriticText (no-failure-case where kinds are named) and the failure-case TestNoTextSaysSevenKinds all PASS; rhiza ship generate exit 0 (0 files changed), rhiza render exit 0. Diff touches only internal/guide text, guide.go's plan-reviewer description, the three plan review tests and the builder's actor record; run.Kinds and severities.md untouched; no TODOs. No guide topic or description says seven kinds (grep).

## Links

### Out

- verifies [plan-reviewer-leaf-contract](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviewer-leaf-contract.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
