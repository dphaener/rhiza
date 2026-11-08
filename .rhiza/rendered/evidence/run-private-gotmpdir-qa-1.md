# run-private-gotmpdir-qa-1

evidence · `01M6XQM7KEJ557X2Y8B9JSG6S6` · actor [claude-code-3cc76aaf](../actor/claude-code-3cc76aaf.md)

## result

PASS

## statement

Built rhiza at the branch head (go build -trimpath) and ran it by hand in a scratch root/repo with RHIZA_* unset. rhiza run start --foreground on a command-step workflow: the step saw GOTMPDIR=<repo>/.git/rhiza/runs/<run id>/gotmp, the folder existed when the step ran (ls -ld), and a real `go build -x` inside the step reported WORK=<that gotmp>/go-build303976419, so go's $WORK sat outside the shared TMPDIR. After the run succeeded, and after a second run failed, runs/ held only the .jsonl journals with no <run id>/ folder left. A third run whose step ran kill -9 $PPID (rhiza died, exit 137) left runs/<id>/gotmp behind; rhiza run ls listed it as a run of its journal, not of the folder; rhiza run abandon <id> exited 0 and the <id>/ folder was gone. Agent steps (session_run GOTMPDIR) were not run by hand, since that would start real agents; they are covered only by the unit tests. go test -trimpath ./internal/run -run GoTmp|TmpDir: ok.

## Links

### Out

- verifies [run-private-gotmpdir](../capability/rhiza/workflows/workflow-runs/run-private-gotmpdir.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
