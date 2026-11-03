# daemon-headless-run-drivers-verified-2

evidence · `01M6GFGGZ0KYFNCQD101Q0CY47` · actor [claude-code-fcf4e672](../actor/claude-code-fcf4e672.md)

## result

PASS

## statement

Verifier claude-code-fcf4e672, attempt 2, at 2d1f0db5. go test ./cmd/rhiza -count=1 -run TestScripts/(daemon-headless-runs|plan-requests|daemon-resume|daemon-picks|run-start-queue)$ passes, and daemon-headless-runs passed 5 more runs in a row. daemon-headless-runs runs rhiza daemon --once against a fake Rex. Rex gets 0 'window new' calls. The child journals its run, rhiza run ls --json shows it live after the daemon has exited, and the script polls the run's .git/rhiza/logs/<cs>.log until the child's own line 'rhiza: step build/build: entered' appears, which answers review finding f1. daemon-picks shows the detached 'rhiza run start --foreground <cs> --permission-mode manual' command, with --autonomous only under --once --autonomous. run-start-queue keeps the queue count at 1 after the started child is live, so the child adds no queue entry. plan-requests shows the 'run start --foreground --workflow plan --new <intent> --request <id>' command, with --milestone for a milestone request and --autonomous only with plan autonomy auto. daemon-resume shows 'rhiza run resume --foreground INTERRUPTED1'. In the diff, detachProgram (cmd/rhiza/detach_unix.go) now runs go cmd.Wait() instead of Process.Release, so the daemon reaps an exited driver instead of leaving a zombie that PidAlive would report live (quality finding f1). make fmt-check, go vet ./... and make self-check pass. make check's test stage failed only on read-timing budgets (internal/api TestReadBudget, internal/slice TestBudget, internal/tui/modelapp TestGeneratedModel, which logged other packages running beside it) at load average about 20; the commit changes no file under internal/, so those failures come from load, not from this leaf.

## Links

### Out

- verifies [daemon-headless-run-drivers](../capability/rhiza/orchestration/daemon/daemon-headless-run-drivers.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
