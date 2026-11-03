# run-continue-through-daemon-verified-2

evidence · `01M6GJ2VT3AHHNQT7ESDK0KKZN` · actor [claude-code-4d7b2408](../actor/claude-code-4d7b2408.md)

## result

PASS

## statement

Re-verified at c9470e92 (its own diff is model-only: the inbox-app-actions amendment to R queueing run.resume, plus actor records). go test ./cmd/rhiza -run 'TestScripts/(run-continue-queue|run-resume.*|daemon-resume|tui-inbox-app|inbox-signoff)$' -count=1: ok. run-continue-queue.txtar shows request-changes - (stdin), approve and resume each journal the action, add the run to .git/rhiza/continues.json and print 'Queued the continuation ... waits for one'; a second action on a queued run is refused (exit 2, nothing journalled); rhiza daemon --max-runs 1 --once resumes one run detached, rhiza run ls --json then shows it live and running after the daemon exits, the log file exists, and rex.log has no 'window new ... run resume'; --foreground --dry-run drives in place with nothing queued. go test ./internal/tui/inboxapp -count=1: ok; TestApproveAndRequestChangesCallTheOperations and TestRetryCallsResume show y/c/R calling run.approve, run.request-changes and run.resume (<run> --retry on a failed run) through the fake runner only; the package has no exec or rex calls. make check: fmt-check, vet and tests passed; a later full go test ./... had one daemon-status failure under load that passed 3/3 on its own (-count=3), and rhiza check now shows no failures.

## Links

### Out

- verifies [run-continue-through-daemon](../capability/rhiza/orchestration/daemon/run-continue-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
