# runs-sessions-apps-verified

evidence · `01M641J09W5W3X0WZX5EE9Y6DA` · actor [claude-code-807c7e32](../actor/claude-code-807c7e32.md)

## result

PASS

## statement

Verifier at 9f5ffb2: go test ./cmd/rhiza -run TestScripts/tui-runs-sessions-apps passes; on its fixture (a passes-workflow run whose build step succeeded, one agent session) `rhiza tui app run-status --once` and `rhiza tui app session-ls --once` are cmp-identical to `rhiza run status --json` and `rhiza session ls --json`, the .rhiza/log file list is identical before and after, and no daemon.lock or socket exists. internal/tui/readapps unit tests pass (both apps register; Op maps run-status to run.status and session-ls to session.ls). A scratch build on this repo's own .rhiza also gave byte-identical output for both apps against the commands. Read the diff: the apps call the api operation in-process with record=false and dryRun unchanged, the same path as the other read apps. make check exits 0.

## Links

### Out

- verifies [runs-sessions-apps](../capability/rhiza/runs-sessions-apps.md)

## History

- `01M62JMTT1Z8FP0SC619V3QS68` Read-only Rex apps for run status and session ls
