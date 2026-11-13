# brief-files-and-invoke-placement

decision · `01M7AHP65F98YYYZH4APSDHPZR`

## reasoning

Read the question, this change set's slice (run-brief-files, agent-invoke-template, guide-contract) and the code. run-brief-files already describes handback.md's content (the cause from rhiza ship check --json), so the format belongs here and the sibling only triggers; one writer avoids two formats. brief.md's sections follow the run-brief-files statement and its acceptance (TestBriefFiles names change set, intent, leaves with acceptance, decisions, constraints, prePushCheck, contract section), kept deterministic per the same-model-same-output constraint. internal/run/store.go keeps journals at <commondir>/rhiza/runs/<id>.jsonl and its listing skips directories, so <id>/ beside it is safe; TmpDir is removed at run end, so it cannot hold files a hand-back or later attempt reads, and per-session subfolders avoid clashes between parallel builder sessions of one leaves step. internal/workflow/workflow.go has prompt on agent steps, Role (builder, verifier, critic, fixer) and fileCheck, so invoke as prompt's alias must apply in each of those places.

## statement

handback.md: this change set owns the file's path, format and writer (run-brief-files ships a function that writes handback.md from a hand-back cause of any failing step: the step and outcome, a gate step's failing commands with their kept output, rhiza ship check --json's failed CI runs with their saved logs, or conflicted files; tested directly); the sibling's run-handback leaf owns triggering a hand-back and calls that writer, writing no handback.md format of its own. brief.md: deterministic Markdown, same bytes for the same model and run, with sections in this order: the change set (title, ID, intent); the leaves (handle, statement, acceptance each); decisions (statement each); constraints (statement each); check commands (rhiza.ship.prePushCheck and rhiza.ship.localGate when set, else none); the contract (a short summary of the obligations ending 'Full contract: rhiza guide contract'). Location: a durable per-run directory beside the journal, <commondir>/rhiza/runs/<run id>/, with one subdirectory per session, <step>-<visit>/<session>, <session> being agent for an agent step, <role>-<leaf>-<attempt> for a leaves step's builder or verifier, <check>-<role>-<cycle> for a checks step's check or critic, <role>-<leaf>-<cycle>-<attempt> for a checks step's leaf fixer or verifier and fixer-<cycle>-<attempt> for its change set fixer (e.g. build-1/builder-<leaf>-1), so parallel sessions and later visits or cycles never share one; it holds that session's brief.md, slice.json and handback.md, is kept when the run ends (unlike the run's temp folder) and is not inside any worktree. invoke and handback keys: accepted everywhere prompt is accepted today: an agent step, a leaves step's builder and verifier, a checks step's each check, critic, fixer and verifier; prompt stays an alias of invoke there, and setting both prompt and invoke on one table is refused by rhiza workflow check.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-e520f785-1](../question/question-e520f785-1.md)
- derived-from [every-agent-session-plan-interview](../source/every-agent-session-plan-interview.md)

### In

- decided-by [run-brief-files](../capability/rhiza/workflows/workflow-runs/run-brief-files.md)
- decided-by [session-brief-wiring](../capability/rhiza/workflows/workflow-runs/session-brief-wiring.md)
- decided-by [agent-invoke-template](../capability/rhiza/workflows/workflow-definitions/agent-invoke-template.md)
- decided-by [arbiter-advisor-brief](../capability/rhiza/workflows/workflow-runs/arbiter-advisor-brief.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
