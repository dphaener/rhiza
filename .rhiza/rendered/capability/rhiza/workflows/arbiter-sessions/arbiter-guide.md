# arbiter-guide

capability · `01M5WRHFC1A3WK2A6G4VQMPTM6` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

rhiza guide ls lists the arbiter topic, rhiza guide arbiter holds an apply recipe for the answering Decision with its wording and reasoning fields and the answers link, and the guide tests pass with the topic within its size cap.

## statement

rhiza guide arbiter tells an arbiter session what to do with the question ID it is given: read the question, its slice and why, decide, write one Decision that answers the question with wording (assumed, would have escalated) and a reasoning field, then exit, writing nothing else.

## Links

### Out

- decided-by [arbiter](../../../../decision/arbiter.md)

### In

- depends-on [arbiter-step](arbiter-step.md)
- verifies [arbiter-guide-verified](../../../../evidence/arbiter-guide-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
