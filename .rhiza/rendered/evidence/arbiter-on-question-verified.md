# arbiter-on-question-verified

evidence · `01M5XD00CZ04HD7HXHXS2BEC7P` · actor [claude-code-54c93d5c-verifier](../actor/claude-code-54c93d5c-verifier.md)

## result

PASS on 2026-10-26 at commit 9b213db (code as of 8b72ad5) on branch cs/autonomous-mode-and-the, re-run with -count=1 by an independent verifier (claude-code-54c93d5c-verifier)

## statement

Testscript arbiter-question with a fake agent that reports a blocked question through the real hook: autonomous, the run's change set gains one question node whose statement is 'Which colour?' with the asking session, the journal has an asked entry and an arbiter entry naming its ID, and the fake arbiter's prompt carries that ID; a fake arbiter that writes nothing ends the run failed with 'arbiter failed on question question-<sid8>-1 (<id>): no live Decision answers question'; normal, run status shows blocked with '- Question: Which colour?', the question count stays 2 and no further arbiter is launched, and the owner is told of the block. Engine tests TestAutonomousQuestionGoesToArbiter, TestAutonomousArbiterFailFailsRun, TestNormalRunQuestionNotArbitrated, TestAutonomousOtherBlockNotArbitrated, TestAutonomousResumeRejoinsArbiter pass.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/arbiter-question$'; go test -count=1 ./internal/run -run 'Autonomous|Normal'

## Links

### Out

- verifies [arbiter-on-question](../capability/rhiza/workflows/arbiter-sessions/arbiter-on-question.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
