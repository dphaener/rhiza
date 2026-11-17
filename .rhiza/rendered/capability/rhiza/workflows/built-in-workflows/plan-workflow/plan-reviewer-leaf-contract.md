# plan-reviewer-leaf-contract

capability · `01M7KJ4TDSR009HWJED70174PC` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a plan-reviewer test finding eight kinds with no-failure-case after oversized-leaf, the narrowed missing-acceptance wording, and a Must not line naming rhiza run review-check; the reviser and critic tests finding no-failure-case wherever they list kinds; and (failure case) no topic or description still saying seven kinds.

## statement

`rhiza guide plan-reviewer` lists eight kinds, adding `no-failure-case` (an acceptance whose failure case is not a real failure path of the leaf, or a failure_waiver with no real reason); narrows `missing-acceptance` to an acceptance that is vague, could not run or runs the full suite; and says never to record what `rhiza run review-check` already checks (a leaf missing statement, acceptance, files or must_not, or with no failure case words and no waiver) or what `rhiza check` reports, such as duplicate handles; `rhiza guide plan-reviser` and `plan-critic` name no-failure-case wherever they list the kinds. Files: internal/guide/text/plan-reviewer.md, internal/guide/text/plan-reviser.md, internal/guide/text/plan-critic.md, internal/guide/guide.go (the descriptions naming seven kinds) and internal/guide's plan review tests. Must not: change run.Kinds, severities.md or any Go outside internal/guide.

## Links

### Out

- depends-on [no-failure-case-kind](no-failure-case-kind.md)
- depends-on [leaf-contract-review-check](leaf-contract-review-check.md)
- decided-by [no-failure-case-major](../../../../../decision/no-failure-case-major.md)

### In

- verifies [plan-reviewer-leaf-contract-verified](../../../../../evidence/plan-reviewer-leaf-contract-verified.md)
- verifies [plan-reviewer-leaf-contract-qa-1](../../../../../evidence/plan-reviewer-leaf-contract-qa-1.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
