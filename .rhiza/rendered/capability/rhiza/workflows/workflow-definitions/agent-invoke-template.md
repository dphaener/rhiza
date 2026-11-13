# agent-invoke-template

capability · `01M79KPBF6JT17PFXC2KD2DW6F` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

go test -trimpath ./internal/workflow -run TestInvokeTemplate covers invoke and handback parsing on an agent step, a leaves step's builder and verifier, and a checks step's check, critic, fixer and verifier, prompt as an alias of invoke, and rendering every placeholder; failure cases: rhiza workflow check refuses an unknown placeholder naming it, and a table setting both prompt and invoke; go test -trimpath ./internal/api -run TestInvokePrompt shows a builder at attempt 2 starting with its rendered invoke followed by the attempt line, and a leaf fixer and a change set fixer starting with their rendered invoke followed by their confirmed findings and answer lines, as with prompt today; go test -trimpath ./cmd/rhiza -run TestScripts/workflow-invoke shows a fake agent receiving '/myskill:conduct --prompt <path to its brief.md>' as its first prompt. Touches: internal/workflow/workflow.go, internal/workflow/check.go, internal/workflow/spec.go, internal/workflow/include.go, internal/workflow/workflow_test.go, internal/api/run.go, internal/api/run_leaves.go, internal/api/run_checks.go (the render calls), cmd/rhiza/testdata/script/workflow-invoke.txtar (new). Must not send hand-backs; run-handback does.

## statement

Wherever a workflow file accepts prompt today (an agent step, a leaves step's builder and verifier, a checks step's check, critic, fixer and verifier) it also accepts invoke, the template of the text the session starts with, and handback, the text a hand-back sends it, with prompt an alias of invoke; both take the placeholders {brief}, {slice} and {handback} (the session's file paths), {change}, {run}, {leaf}, {attempt} and {worktree} beside today's; the daemon renders invoke and, as it does with prompt today, appends after it unchanged the lines it adds for that session (a retried or freshly-rounded builder's retry, Decision-answer and attempt lines from leafPrompt, a leaf fixer's and a change set fixer's confirmed findings and answer lines from fixerPrompt and changeSetFixerPrompt), and the agent's adapter delivers the result as the first prompt; with no such lines the first prompt is exactly the rendered invoke; handback is rendered for run-handback to send.

## Links

### Out

- depends-on [session-brief-wiring](../workflow-runs/session-brief-wiring.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)
- decided-by [brief-files-and-invoke-placement](../../../../decision/brief-files-and-invoke-placement.md)

### In

- verifies [agent-invoke-template-verified](../../../../evidence/agent-invoke-template-verified.md)
- verifies [agent-invoke-template-qa-1](../../../../evidence/agent-invoke-template-qa-1.md)
- verifies [agent-invoke-template-verified-2](../../../../evidence/agent-invoke-template-verified-2.md)
- verifies [agent-invoke-template-qa-2](../../../../evidence/agent-invoke-template-qa-2.md)
- verifies [pr-565-every-agent-session-a-run](../../../../evidence/pr-565-every-agent-session-a-run.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
