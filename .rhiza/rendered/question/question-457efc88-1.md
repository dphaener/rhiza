# question-457efc88-1

question · `01M6WGRE5F081ZYEVZHNMY0DZ3`

## session

457efc88-a06a-4e5e-afa5-01a1b907c681

## statement

Fix A Go build inside a run serves no intent yet. The last bug fix (Fix Review question handles) served milestone-29 directly; milestone 30 is about sharing Rhiza publicly. Which should this change set serve? / No code was found that deletes a fresh ~/tmp/go-build* folder: gosweep needs a 2h-old top-level mtime (this WORK lived 32s, 00:20:00-00:20:32Z), and no cron, timer, tmpfiles rule or agent transcript removed one. How should the change set stop it? / Where should the retry on go's vanished-$WORK failure live?

## Links

### In

- answers [go-work-vanish-fix-shape](../decision/go-work-vanish-fix-shape.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
