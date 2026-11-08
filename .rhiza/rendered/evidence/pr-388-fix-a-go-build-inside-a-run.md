# pr-388-fix-a-go-build-inside-a-run

evidence · `01M6XRK55VVSC6F1DTKQBCPZK7` · actor [dphaener](../actor/dphaener.md)

## head

71ed9b211de6e65d2e63ccd59d6d9c2d4faabaf0

## pr

https://github.com/dphaener/rhiza/pull/388

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #388 (https://github.com/dphaener/rhiza/pull/388) passed CI at 71ed9b2

## status

verified

## Links

### Out

- verifies [run-private-gotmpdir](../capability/rhiza/workflows/workflow-runs/run-private-gotmpdir.md)
- verifies [step-retry-vanished-go-work](../capability/rhiza/workflows/workflow-runs/step-retry-vanished-go-work.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
