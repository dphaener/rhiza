# workflow-check-merge-no-gate-qa-1

evidence · `01M7E24TT2D99X1YQSZ1BMDKW1` · actor [claude-code-98071bd0](../actor/claude-code-98071bd0.md)

## result

PASS

## statement

QA: rhiza workflow check on a hand-written build -> test -> merge (rhiza ship merge) warns naming merge and the path build -> test -> merge, in text and --json, exit 0 and valid true; a merge (multi-line run with gh pr merge) reached only via a kind=gate step: no warning; echo "rhiza ship merged" is not taken for a merge; ship-change: none. Defect: built-in plan, plan-milestone, plan-review and ship-change-classic, whose merge follows a command step running rhiza run gate or rhiza ship gate, get the warning (finding recorded).

## Links

### Out

- verifies [workflow-check-merge-no-gate](../capability/rhiza/workflows/workflow-authoring/workflow-check-merge-no-gate.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
