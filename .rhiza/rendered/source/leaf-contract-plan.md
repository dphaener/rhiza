# leaf-contract-plan

source · `01M7KJ3AWBTJ2ZYPF20HGSZ2MA`

## quotes

```json
[]
```

## raw_sha256

af50ab966583d458ac701861664e32d89be2f519fa50181ee3cacde652415498

## session

d69efb4e-12a2-4bd8-8dc8-01a1e71e989c

## summary

Autonomous conductor run 01M7KHWS planned change set 01M76XT3 (Leaf contract) with no owner. Research: prompt-gap.md rec 7 and section 2.1 items 2-5 ask for per-leaf files and must-not, a failure case or a waiver in every acceptance, one-concern sizing, a no-failure-case plan-review kind, and mechanical checks moved from the reviewer to review-check; section 5 names Go model fields files and must-not, interim inside acceptance text. Code read: internal/op/fields.go refuses unknown fields on new writes (knownFieldsOnly in internal/api/write_node.go); run.Kinds in internal/run/journal.go holds seven kinds, uncovered-intent and oversized-leaf already landed with Plan role topics (01M76XRP), so this change set adds only no-failure-case (the scope note 01M76YY); review-check's plan-review path (runRunReviewCheck in internal/api/run_review.go) reads rounds only and checks no leaf field; brief.md (internal/run/brief.go, internal/api/run_brief.go) lists each leaf's statement and acceptance; landed acceptances already mark a failure case with the words (failure case); rhiza check already reports duplicate handles. Planner rulings, no owner asked: the contract is capability fields files, must_not and failure_waiver rather than acceptance text; review-check tests the fields and the words failure case at every round; no-failure-case is major and judges whether the failure case is real. This change set's own leaves carry files and must-not in their statement text, because the fields do not exist until its first leaf is built.

## title

Leaf contract plan (autonomous conductor)

## Links

### In

- derived-from [leaf-contract-as-fields](../decision/leaf-contract-as-fields.md)
- derived-from [failure-case-words](../decision/failure-case-words.md)
- derived-from [leaf-contract-check-every-round](../decision/leaf-contract-check-every-round.md)
- derived-from [no-failure-case-major](../decision/no-failure-case-major.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
