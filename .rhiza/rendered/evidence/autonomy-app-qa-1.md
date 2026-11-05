# autonomy-app-qa-1

evidence · `01M6DMT3KRWWVGZ4ZE28T4T1D0` · actor [claude-code-f364d51f](../actor/claude-code-f364d51f.md)

## result

PASS

## statement

At 69611388, built rhiza to a scratch path and ran it in a scratch repo. `rhiza tui app autonomy --once` cmp-equals `rhiza daemon settings --json`. Ran the live app in a Rex window and drove it with rex send-key: enter on rhiza.daemon.paused set it true (source git config); enter on rhiza.autonomy.plan toggled owner->auto and back; editing maxRuns to 0 was refused inline ('must be a whole number of at least 1') with the entry kept open, then 5 saved; permissionMode edited as a word. d then n did nothing; d then y opened a Rex tab named 'rhiza daemon' and, with permissionMode acceptEdits, the header became 'Daemon running, process Mac.localdomain:79345'; d then y again sent TERM and the process was gone, header 'Daemon not running'. With permissionMode unset (the default) or 'default', d y still printed 'Started the daemon' while `rhiza daemon --autonomous` exited at once with a permission-mode error and its tab closed (reported as a finding). go test ./internal/tui/autonomyapp and TestScripts/(tui-autonomy-app|rex-install|describe-golden) pass.

## Links

### Out

- verifies autonomy-app (retired)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
