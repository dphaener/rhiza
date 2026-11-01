# daemon-status-version-qa-1

evidence · `01M6AQSE7NVRW1GM4RWAJQTF6K` · actor [claude-code-17282cde](../actor/claude-code-17282cde.md)

## result

PASS

## statement

At 8a20d26a, with a binary built with -X main.version=qa-v1 running `rhiza daemon` in a scratch repo, `rhiza daemon status` printed 'Version: qa-v1' and --json carried "version": "qa-v1". After rebuilding in place as qa-v2 and the daemon re-executing, status printed 'Version: qa-v2'. After SIGINT, status printed 'Daemon not running', 'Last daemon: ... stopped ...' and still 'Version: qa-v2'. go test ./cmd/rhiza -run 'TestScripts/(daemon-status|describe-golden)$' passes.

## Links

### Out

- verifies [daemon-status-version](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-version.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
