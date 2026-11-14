# workflow-check-uncapped-loop-qa-1

evidence · `01M7E24TXX0RJYKFZZY9Y4D533` · actor [claude-code-98071bd0](../actor/claude-code-98071bd0.md)

## result

PASS

## statement

QA: test on.fail = "build" back to agent step build warns naming step, outcome, target and handback:build, in text and --json, exit 0; on.fail = "handback:build": no warning; ship-change and a fresh rhiza workflow new: no loop warning (only the placeholder one). Rough edge: merge on.fail = "test" where test routes on.fail = "handback:build" still warns although the hand-back is counted (minor finding).

## Links

### Out

- verifies [workflow-check-uncapped-loop](../capability/rhiza/workflows/workflow-authoring/workflow-check-uncapped-loop.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
