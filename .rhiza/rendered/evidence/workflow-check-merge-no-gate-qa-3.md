# workflow-check-merge-no-gate-qa-3

evidence · `01M7E42S7F4NEQ5X9YQAVWJGCM` · actor [claude-code-2248efc5](../actor/claude-code-2248efc5.md)

## result

PASS

## statement

QA cycle 3: no built-in (plan, plan-milestone, plan-review, ship-change-classic, ship-change) warns of an ungated merge; build -> gate(echo rhiza ship gate) -> merge now warns naming the path (earlier qa f4 fixed); a multi-line 'cd . &&\n rhiza run gate' gates; rhiza ship gate --failed does not gate; all exit 0, valid true. TestScripts workflow-check-merge passes.

## Links

### Out

- verifies [workflow-check-merge-no-gate](../capability/rhiza/workflows/workflow-authoring/workflow-check-merge-no-gate.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
