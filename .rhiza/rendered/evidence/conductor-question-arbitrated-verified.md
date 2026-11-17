# conductor-question-arbitrated-verified

evidence · `01M7M7QAM5GF2R4THVNS4CPTG6` · actor [claude-code-496ebf13](../actor/claude-code-496ebf13.md)

## result

PASS

## statement

Detached checkout at d0543beab. Acceptance: go test -p 2 -trimpath ./internal/run -run 'TestConductorQuestion|TestNormalRunQuestionNotArbitrated|TestAutonomousQuestionGoesToArbiter|TestAutonomousArbiterFailFailsRun' exit 0; TestConductorQuestionArbitratedInNormalRun passes for included=false (conduct) and included=true (ship/conduct via a workflow step) in a non-autonomous run, asserting one asked entry, the arbiter on Q1, the answer delivered, no advisor entry or StartAdvisor call and no notices, plus the failure case (arbiter fails -> step error naming question Q1); TestNormalRunQuestionNotArbitrated still passes (owner notices, no arbiter or asked entry). Checks block re-run: go test -p 2 -trimpath ./internal/run/... exit 0; go vet ./internal/run exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 and the tree stayed clean. Diff read: arbitrates(step, b) and arbitratesIn(step) = Autonomous || def.Steps[step].Conducts() keyed in block, unblock, pollAgent, adviseQuestion, arbitrateAsked and the check/leaf/review-fix polls (those are review steps, so their behaviour is unchanged); still question-only; State.Autonomous, signoff, gates and daemonStart untouched; no stubs or TODOs.

## Links

### Out

- verifies [conductor-question-arbitrated](../capability/rhiza/workflows/arbiter-sessions/conductor-question-arbitrated.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
