# planner-leaf-contract-verified

evidence · `01M7KS9RJC10325DMHTEP2JF2V` · actor [claude-code-6d54606b](../actor/claude-code-6d54606b.md)

## result

PASS

## statement

In a detached checkout at 364011c48 re-ran the Checks: block: go test -trimpath -p 2 ./internal/guide/... exit 0 (TestPlanLeafContract, TestCaps PASS with -count=1 -v), rhiza ship generate exit 0 with 0 changed files, rhiza render exit 0 with a clean tree. Diff read: plan-recipes Leaves names files, must_not, failure_waiver, `failure case`, one fresh session and the split rule, and its apply recipe carries files= and must_not= with a (failure case) acceptance; planner Expected outcome and self-check name files, must_not and the failure case; the test checks each, plus both topics within TopicCap (failure case; planner 860 tokens, plan-recipes about 882, cap 1000). Only plan-recipes.md, planner.md, the plan-recipes description in guide.go, plan_test.go and the builder's actor record changed; no plan reviewer, critic or reviser topic touched; no stubs or TODOs.

## Links

### Out

- verifies [planner-leaf-contract](../capability/rhiza/workflows/built-in-workflows/plan-workflow/planner-leaf-contract.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
