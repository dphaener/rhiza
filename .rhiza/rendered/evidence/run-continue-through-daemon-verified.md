# run-continue-through-daemon-verified

evidence · `01M6G3JCGR7GPQ138AR92GNBB6` · actor [claude-code-d272bcd9](../actor/claude-code-d272bcd9.md)

## result

PASS

## statement

At 8135c5f1 (leaf commits f805ea0e, 8135c5f1) I ran go test ./cmd/rhiza -count=1 -run 'TestScripts/(run-continue-queue|review-signoff|run-resume|run-resume-retry|inbox-signoff|describe-golden|run-abandon-failed)$' (ok), and run-continue-queue three more times (ok each). In run-continue-queue, rhiza run request-changes - <run> (findings on stdin), rhiza run approve and rhiza run resume each journal the action (changes/approved sign-off), add the run to .git/rhiza/continues.json and print 'Queued the continuation ... no daemon runs, so it waits for one'. A second continuation is refused (exit 2, nothing journalled). rhiza daemon --max-runs 1 --once continues only the oldest run, as a detached rhiza run resume --foreground child with output in .git/rhiza/logs/<run>.log and no 'window new ... run resume' in the fake rex log. After the daemon exits, rhiza run ls --json shows that run live:true, state running (its check agent is held). Released, the child drives review visit 2. A second daemon --once then takes the approved run to succeeded. --foreground drives in place for resume (run-continue-queue, run-resume, run-abandon-failed) and for approve and request-changes (review-signoff). go test ./internal/tui/inboxapp -count=1 passes: TestApproveAndRequestChangesCallTheOperations and TestRetryCallsResume show the y, c and R keys calling run.approve, run.request-changes (findings file, removed afterwards) and run.resume [--retry] through the fake runner. The inbox model no longer has an Exec path, so it cannot run rex window new. go test ./internal/api -run 'Run|Daemon|Signoff|SignOff|Inbox|Continue|Resume|Approve|Request' and ./internal/daemon pass. A full ./internal/api run timed out in TestReadBudget and then TestLinkEndpointKinds, model-read timing tests outside this leaf, at load average ~20 from other sessions. I read daemonQueuedContinues: it starts queued runs oldest first, within max, through daemonDetach, and leaves them out of the daemon's own resume. Both attempt-1 findings (request-changes - not covered; run ls live after the daemon exits not shown) are fixed.

## Links

### Out

- verifies [run-continue-through-daemon](../capability/rhiza/orchestration/daemon/run-continue-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
