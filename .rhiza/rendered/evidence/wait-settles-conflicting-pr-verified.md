# wait-settles-conflicting-pr-verified

evidence · `01M6AJ6NTXHM2PEJVT9N3HNHRV` · actor [claude-code-2d78a3e9](../actor/claude-code-2d78a3e9.md)

## result

PASS

## statement

At 2cfda3de, go test ./cmd/rhiza -run TestScripts/ship-check-settles-conflict passes: its fake gh reports mergeable CONFLICTING until origin's branch holds main; with a conflict only in .rhiza/rendered, ship check --settled exits non-zero printing 'rebased onto main and pushed', origin/cs/ship-it then contains origin/main, the ship record's pushed equals the new origin head, and a second --settled check exits 0; with a conflict in code.txt it exits 75, prints a stop: line naming code.txt and rebase-conflict, commits only the question, and HEAD~1 and origin/cs/ship-it equal the pre-check head. Reverting internal/api/ship.go to HEAD~1 makes the script fail (line 50, unexpected success). Diff read: ship.check fetches mergeable, and on --settled + CONFLICTING calls runShipRebase, mapping Conflicts to exit run.StopExit. run-ship-change and describe-golden pass; make check exits 0.

## Links

### Out

- verifies [wait-settles-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-settles-conflicting-pr.md)

## History

- `01M6A0JJDJ0FXD58K4WV88FV0A` A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict
