# workflow-check-uncapped-loop-qa-2

evidence · `01M7E3FQYRA9NSWDJKJWADNCW9` · actor [claude-code-a8f5980f](../actor/claude-code-a8f5980f.md)

## result

PASS

## statement

QA cycle 2: a test step with on.fail = build warns in text and --json, exit 0, suggesting handback:build; build -> test(handback:build) -> gate -> merge(on.fail test) no longer warns (earlier qa f3 fixed); ship-change and a workflow from rhiza workflow new raise no loop warning. Rough edge: built-in plan and ship-change-classic raise many loop warnings on max_visits-capped loops (minor).

## Links

### Out

- verifies [workflow-check-uncapped-loop](../capability/rhiza/workflows/workflow-authoring/workflow-check-uncapped-loop.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
