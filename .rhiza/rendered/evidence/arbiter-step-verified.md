# arbiter-step-verified

evidence · `01M5XD004RHXKC40PRFSR37Z9H` · actor [claude-code-54c93d5c-verifier](../actor/claude-code-54c93d5c-verifier.md)

## result

PASS on 2026-10-26 at commit 9b213db (code as of 8b72ad5) on branch cs/autonomous-mode-and-the, re-run with -count=1 by an independent verifier (claude-code-54c93d5c-verifier)

## statement

TestArbiterStepLoads parses a workflow with kind = "arbiter" into the expected Step and Arbiter(); TestArbiterStepRefused refuses prompt ('an arbiter step cannot have prompt: its prompt is fixed') and run, until, until_when, timeout, max_visits, on.*. internal/api TestStartArbiterPrompt, TestArbiterAnswer, TestRunStartChecksArbiter pass. Testscript arbiter-step (landed with leaf arbiter-on-question's commit, since only an autonomous run's question reaches the arbiter) runs the step on a question with a fake adapter: the arbiter's prompt equals exactly 'Decide question <question ID>: follow rhiza guide arbiter.', without the change set ID, a question- handle or the title; a Decision with reasoning passes ('answered by assume-blue'), no Decision fails ('no live Decision answers question'), a Decision without reasoning fails ('has no reasoning'), and the run fails with it.

## tests

go test -count=1 ./internal/workflow -run Arbiter; go test -count=1 ./internal/api -run Arbiter; go test -count=1 ./cmd/rhiza -run 'TestScripts/arbiter-step$'

## Links

### Out

- verifies [arbiter-step](../capability/rhiza/workflows/arbiter-sessions/arbiter-step.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
