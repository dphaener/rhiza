# journal-failure-reported

decision · `01M5XCKYWVNQNA2HYANR8FSEYM`

## statement

If the run's journal cannot be written after an autonomous accept is in the log, the command fails naming the log entry it wrote, and the accept is not rolled back.

## wording

the building agent's choice while building leaf autonomous-gate-rule

## Links

### In

- decided-by [autonomous-gate-rule](../capability/rhiza/workflows/autonomous-runs/autonomous-gate-rule.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
