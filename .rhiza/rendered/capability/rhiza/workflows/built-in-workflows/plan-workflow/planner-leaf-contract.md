# planner-leaf-contract

capability · `01M7KJ4T7B9ARS7P07F9M1G7K2` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a plan-recipes test finding files, must_not, failure_waiver, `failure case` and one fresh session in the Leaves section and its apply recipe JSON carrying files= and must_not= fields; a planner test finding files, must_not and the failure case in Expected outcome and in the self-check step; and (failure case) both topics' expanded bodies within TopicCap.

## statement

`rhiza guide plan-recipes`' Leaves section teaches the leaf contract: its apply recipe writes a leaf with statement, acceptance, files and must_not; each acceptance holds at least one failure case marked with the words `failure case`, or the leaf carries failure_waiver with its reason; a leaf is one concern sized for one fresh session, split on two unrelated concerns or a files list too long to hold; and `rhiza guide planner`'s Expected outcome and self-check name files, must_not and the failure case. Files: internal/guide/text/plan-recipes.md, internal/guide/text/planner.md, internal/guide/guide.go (the two descriptions, if they change) and internal/guide's plan tests. Must not: touch the plan reviewer, critic or reviser topics, or push either topic past TopicCap.

## Links

### Out

- depends-on [leaf-contract-fields](leaf-contract-fields.md)
- decided-by [leaf-contract-as-fields](../../../../../decision/leaf-contract-as-fields.md)
- decided-by [failure-case-words](../../../../../decision/failure-case-words.md)

### In

- verifies [planner-leaf-contract-verified](../../../../../evidence/planner-leaf-contract-verified.md)
- verifies [planner-leaf-contract-qa-1](../../../../../evidence/planner-leaf-contract-qa-1.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
