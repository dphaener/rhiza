# ship-command-after-close-out-verified-2

evidence · `01M642M3Q6R66MCA7A8RD8Z3B0` · actor [claude-code-5aef9d8f](../actor/claude-code-5aef9d8f.md)

## result

PASS

## statement

Verifier, attempt 2, at b3f29d4: `go test ./cmd/rhiza -run 'TestScripts/run-ship-change$' -count=1` passes. The script runs `rhiza run approve` from the change set's worktree $WORK/repo.ship-it. Close-out removes that folder, the ship step then runs the repo's .rhiza/workflows/ship.sh in the main checkout, ship.log holds 'shipped' once, approve prints exit=0, and the ship-change run ends succeeded with close-out -> ship -> end all passing. With 1a60493's internal/api/run.go and run_test.go restored in a scratch worktree, the same script fails at `rhiza run approve`. TestStepRunnerDirFallsBackToMainCheckout passes: with the worktree gone and git failing, dir() gives the main checkout joined with the project subfolder. go test ./... passes and vet and gofmt are clean. make self-check fails only on stale .rhiza/rendered pages for this run's uncommitted session actor records, which the run's render step writes. The fix commit b3f29d4 adds the model links: ship-command-starts-in-main-checkout supersedes ship-command-main-from-common-dir and decides this capability, matching the code (the runner records the main checkout from the common dir at run start). The review's remaining minor points still stand: the script does not record pwd, and a hand-run `rhiza ship command` from a removed folder still does nothing. Both fall outside the acceptance as decided.

## Links

### Out

- verifies [ship-command-after-close-out](../capability/rhiza/ship-command-after-close-out.md)

## History

- `01M625RC3X2VA17M5V2HVS3HBP` ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.
