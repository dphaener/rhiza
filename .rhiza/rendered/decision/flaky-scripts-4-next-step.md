# flaky-scripts-4-next-step

decision · `01M78F0DNBQMF216BY87HNJRWP`

## reasoning

The question gives no context beyond 'What would you type next?', so it reads as a session asking how to proceed. The change set's intent is to reproduce the UNKNOWN STEP flake from PR #498 and then fix or quarantine it. The CI logs for that head show the step was 'Run scripts' in scripts (4/4) and that the failing test was inbox-stalled.txtar:32. AGENTS.md says to reproduce a flake locally and repeatedly by injecting a delay at the race (no stress runs on the shared host), to fix it or quarantine it with a bug filed, and to push only after it passes locally. So the answer is to start that procedure on inbox-stalled.

## statement

Next step: reproduce the scripts (4/4) flake locally. CI run 43256685725 attempt 1 (PR #498, head d1e350c) failed in TestScripts/inbox-stalled at testdata/script/inbox-stalled.txtar:32 (unexpected match for "stalled" in stdout), and attempt 2 passed. Run `go test -trimpath ./cmd/rhiza -run TestScripts/inbox-stalled` several times, inject a delay at the suspected race (the stalled-detection clock or timer), then fix the timing dependence or quarantine the script with a bug filed through `rhiza request --kind bug`. Push only after it and the touched packages pass locally.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-aba0b284-1](../question/question-aba0b284-1.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
