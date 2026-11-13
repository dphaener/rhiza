# agent-invoke-template-verified

evidence · `01M7AW29B7SE2DVYMSD4E39CF3` · actor [claude-code-751e52ad](../actor/claude-code-751e52ad.md)

## result

PASS

## statement

At 843b0e128 (attempt 2, over 1e0cc2dcf) re-ran the Checks block and acceptance: go test -trimpath ./internal/workflow/... ok; ./internal/workflow -run TestInvokeTemplate PASS (incl. unknown invoke/handback placeholder, prompt+invoke both set, agent-step {leaf}/{attempt} rendering empty); ./internal/api -run TestInvokePrompt PASS (builder attempt 2 = rendered invoke then attempt line; leaf and change set fixers = rendered invoke then confirmed findings and answer lines); ./cmd/rhiza TestScripts/(workflow-invoke|workflow-check|workflow-describe|docs-reference-workflow|describe-golden) ok, workflow-invoke shows '/myskill:conduct --prompt <brief.md path>'. Scratch-built rhiza workflow check: agent step invoke/handback 'do {leaf} {attempt}' valid (exit 0), resolving attempt-1 finding f1; {bogus} and {check} refused naming the placeholder (exit 1); prompt+invoke refused (exit 1). rhiza ship generate exit 0 (0 changed), rhiza render exit 0, tree clean. Diff read: render calls in run.go/run_leaves.go/run_checks.go use session file paths via promptVars; no hand-back sending, no TODOs or stubs.

## Links

### Out

- verifies [agent-invoke-template](../capability/rhiza/workflows/workflow-definitions/agent-invoke-template.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
