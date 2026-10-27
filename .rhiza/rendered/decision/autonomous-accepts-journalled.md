# autonomous-accepts-journalled

decision · `01M5XCKVV63JPRZWRRZZX5H26W`

## statement

An accept under rule autonomous also writes an accepted entry, naming the change set and the log entry, to the journal of the run it was made in, and run status reads what the run accepted from that journal.

## wording

the planning agent's option, picked by the orchestrator acting as arbiter for the owner

## Links

### In

- decided-by [autonomous-gate-rule](../capability/rhiza/workflows/autonomous-runs/autonomous-gate-rule.md)
- decided-by [autonomous-review](../capability/rhiza/workflows/autonomous-runs/autonomous-review.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
