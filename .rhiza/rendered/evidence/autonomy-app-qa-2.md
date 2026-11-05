# autonomy-app-qa-2

evidence · `01M6DPJQ2HAJXJ1P45205MKRJZ` · actor [claude-code-b3f5c4d8](../actor/claude-code-b3f5c4d8.md)

## result

PASS

## statement

At 7f653b81 (after fix 1a5e3269), built rhiza to a scratch path and ran `rhiza tui app autonomy` live in a Rex window over a fresh scratch repo, driven with rex send-key. With permissionMode unset, d then y said 'Starting the daemon in the Rex tab rhiza daemon', then within one reload 'The daemon did not start: its error is in the Rex tab rhiza daemon'; that tab stayed open showing the daemon's 'no permission mode' error (the cycle-1 QA finding is fixed). Edited permissionMode to acceptEdits in place (source became git config); d y then showed 'The daemon is running in the Rex tab rhiza daemon' with header 'Daemon running, process Mac.localdomain:54983'; d y again sent TERM, the process was gone and the header said 'Daemon not running'. Enter toggled paused false->true and plan owner->auto, both from git config. go test autonomyapp and TestScripts/(tui-autonomy-app|rex-install) pass.

## Links

### Out

- verifies autonomy-app (retired)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
