# agent-invoke-template-qa-1

evidence · `01M7AWMNY1N8Z5G63F673QHCEM` · actor [claude-code-1707fde5](../actor/claude-code-1707fde5.md)

## result

PASS

## statement

QA with a scratch build (go build -trimpath) in hermetic repos with fake rex/claude: agent step invoke rendering every placeholder ({leaf}/{attempt} empty) delivered as the first prompt, brief/slice/handback paths under .git/rhiza/runs/<run>/build-1/agent/; prompt alias works; leaves builder attempt 2 = rendered invoke then the attempt line; leaf and change set fixers = rendered invoke then confirmed findings lines; check at cycle 1 got exactly its invoke. workflow check refuses unknown placeholders in invoke and handback (agent, verifier, check), prompt+invoke on step, builder, check and critic, handback on command/wait/arbiter and at leaves step level. Rough edges recorded as minor findings f2, f3.

## Links

### Out

- verifies [agent-invoke-template](../capability/rhiza/workflows/workflow-definitions/agent-invoke-template.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
