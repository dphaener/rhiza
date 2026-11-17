# conductor-question-arbitrated-qa-1

evidence · `01M7M99R89GT63EFGCGX140BFH` · actor [claude-code-2166ab54](../actor/claude-code-2166ab54.md)

## result

PASS

## statement

Built ./cmd/rhiza (HEAD 289c5fa77) to a scratch path; ran it in a scratch Rhiza root with a fake claude (real hook path) and fake rex, runs not started --autonomous. Conduct step's own question: asked entry, question node, arbiter pass, answer typed into its tab, run succeeded, no advisor entry and no state=blocked notice; same through workflow outer whose workflow step includes conduct (step conduct/conduct). Failure case: arbiter wrote no Decision, so the run failed with the step's error naming the question ('arbiter failed on question question-...-1 (...)'). Regression: workflow plain blocked on the owner with a state=blocked notice and an advisor entry, no asked entry. Scoped tests: go test ./internal/run -run 'TestConductorQuestion|TestNormalRunQuestionNotArbitrated|TestAutonomousQuestionGoesToArbiter|TestAutonomousArbiterFailFailsRun' ok.

## Links

### Out

- verifies [conductor-question-arbitrated](../capability/rhiza/workflows/arbiter-sessions/conductor-question-arbitrated.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
