# agent-invoke-template-verified-2

evidence · `01M7B1NPF6Q7KCJKYK278GG19Y` · actor [claude-code-1e691e2c](../actor/claude-code-1e691e2c.md)

## result

PASS

## statement

At 32375f2b9 (review cycle 2 fix): go test -trimpath ./internal/workflow/... and ./internal/guide/... exit 0; go test -trimpath ./internal/workflow -run TestInvokeTemplate PASS (agent step, builder/verifier, check/critic/fixer/verifier invoke and handback, prompt alias, every placeholder, {leaf}/{attempt} empty on check and critic, fixer {attempt}; refusals for unknown placeholder naming it and prompt+invoke both set); TestChecksStepRefused PASS; go test -trimpath ./internal/api -run TestInvokePrompt PASS and the commit's scoped api run (TestStartFix|TestFixChangeSet|TestCycleCapFixerPrompts|TestProgramStatusReviewFixers|TestFixLeafEngine|TestOwnerCycleFixerPrompts|TestInvoke|TestBrief|Checks) ok; go test -trimpath ./cmd/rhiza -run TestScripts/workflow-invoke ok. Diff read: FixerVars gains Attempt passed from fixerPrompt and changeSetFixerPrompt, check/critic placeholders gain leaf and attempt, docs/reference/workflow-file.md and contract.md updated to match; no stubs or TODOs. make self-check fails only on this verifier session's own unrendered actor record, not on the commit.

## Links

### Out

- verifies [agent-invoke-template](../capability/rhiza/workflows/workflow-definitions/agent-invoke-template.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
