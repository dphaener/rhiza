# daemon-reexec-new-binary-qa-1

evidence · `01M6AQSCFWACZFPVB0QX4FZ298` · actor [claude-code-17282cde](../actor/claude-code-17282cde.md)

## result

PASS

## statement

At 8a20d26a I built rhiza into a scratch dir with -X main.version=qa-v1, ran `rhiza daemon --permission-mode manual` in a fresh scratch repo (pid 4031), then rebuilt the same path in place as qa-v2. Within one 5s poll stderr printed 'binary changed, re-executing' and 'daemon running, process ...:4031': same pid, argv unchanged (ps shows `rhiza daemon --permission-mode manual`), and daemon status then showed Version: qa-v2. Over the next polls it did not re-exec again. Moving the binary away for two polls (as mid-rebuild) and back with the same mtime caused no exec. After the re-exec a second `rhiza daemon --once` was refused with exit 2 ('another daemon holds the lock: process ...:4031'), so the new image holds the lock. SIGINT then stopped it cleanly ('Daemon ...:4031 stopped'). The pending-start delay and live runs carrying on were not exercised by hand (empty model, no runs); go test ./internal/daemon ./internal/api pass.

## Links

### Out

- verifies [daemon-reexec-new-binary](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
