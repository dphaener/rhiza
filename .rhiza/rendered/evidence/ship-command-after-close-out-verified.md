# ship-command-after-close-out-verified

evidence · `01M6424K6NKST80TB325D0HM89` · actor [claude-code-648c68f4](../actor/claude-code-648c68f4.md)

## result

PASS

## statement

Verified by claude-code-648c68f4 at 349b880. go test ./cmd/rhiza -run TestScripts/run-ship-change passes: the script approves from the change set's worktree $WORK/repo.ship-it, close-out removes that folder under the running process, the ship step (rhiza ship command) runs .rhiza/workflows/ship.sh in the main checkout with exit 0 (ship.log holds 'shipped' once, after main-green), and the run ends 'succeeded' with 'close-out: pass, then ship: exit 0' and 'ship: pass, then end: exit 0'. With internal/api/run.go reverted to its parent in a scratch worktree the same script fails: ship errors with 'chdir $WORK/repo.ship-it: no such file or directory', so the test catches the bug. TestStepRunnerDirFallsBackToMainCheckout passes (main checkout from the git common dir, project subfolder kept, git failing). make fmt-check vet test pass. make self-check fails at this detached commit only because the builder's actor record (claude-code-9c9d7fe0) is still uncommitted in the run worktree, which ship-change's render step commits later; it does not affect the leaf.

## Links

### Out

- verifies [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
