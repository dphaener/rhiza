# milestone-10

intent · `01M6A0KZ2XJBXWDE4ENGH8CD67`

## milestone

10

## statement

Runs heal themselves: a PR that conflicts after main moves is rebased by its own run, a flaky CI check is rerun once before a run fails, a daemon start that never journaled is retried without a restart, the daemon picks up a rebuilt binary by itself, and a session stuck with no activity shows as stalled

## Links

### Out

- decided-by [m10-picked-by-orchestrator](../decision/m10-picked-by-orchestrator.md)

### In

- serves change set `01M6A0JJDJ0FXD58K4WV88FV0A` (A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict)
- serves change set `01M6A0JQ2WHGH7DDFKJGRXJD85` (A failed CI check on a run's PR or main is rerun once per head before the run fails)
- serves change set `01M6A0JVE5ZV75XEMB0ZNYN8SR` (A daemon start that wrote no journal is retried after a hold, without restarting the daemon)
- serves change set `01M6A0K064DH4TP3JJ0861DFG4` (The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version)
- serves change set `01M6A0K4K58YMW04FA38373Y8P` (A session with no activity for 20 minutes shows as stalled in run status and the inbox)

## History

- `01M6A0JENBMBYPYRN5195AGPKH` Plan milestone 10: runs heal themselves, picked by the orchestrator under the owner's standing instruction
