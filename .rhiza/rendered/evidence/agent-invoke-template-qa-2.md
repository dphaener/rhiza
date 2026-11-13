# agent-invoke-template-qa-2

evidence · `01M7B24CRFZAWTR93YSMHB6KQ5` · actor [claude-code-a1cb0102](../actor/claude-code-a1cb0102.md)

## result

PASS

## statement

QA cycle 2, scratch build (go build -trimpath) in a hermetic Rhiza root with fake rex/claude: agent step invoke using every placeholder rendered as the first prompt ({leaf}/{attempt} empty, {brief}/{slice}/{handback} the session's build-1/agent/ paths); leaves builder at attempts 2 and 3 got rendered invoke then the attempt line; checks step check and critic invoke with {leaf}/{attempt} (the cycle-2 fix) pass workflow check and render empty; leaf fixer got rendered invoke with {leaf} and {attempt} then its confirmed-findings lines, change set fixer the same with {leaf} empty and Review: trailer line. workflow check refuses {check} in a fixer and a builder, {cycle} in an agent step, {bogus} in a critic handback, prompt+invoke on a verifier, handback on a command step, each naming it (exit 1). Scoped tests TestInvokeTemplate, TestInvokePrompt, TestScripts/workflow-invoke ok.

## Links

### Out

- verifies [agent-invoke-template](../capability/rhiza/workflows/workflow-definitions/agent-invoke-template.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
