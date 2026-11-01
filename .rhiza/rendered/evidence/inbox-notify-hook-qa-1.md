# inbox-notify-hook-qa-1

evidence · `01M6BQR7RJRTXYA5N4H9B79D7G` · actor [claude-code-a6c7a851](../actor/claude-code-a6c7a851.md)

## result

PASS

## statement

QA at 6f050ab9 with a scratch build (go build ./cmd/rhiza), run by hand outside testscript in a fresh git repo with REX_* and RHIZA_TEST_* unset, over question, signoff and failed run journals from daemon-notify-hook.txtar. rhiza daemon --once --permission-mode manual --notify-command hook.sh ran the hook 3 times: each got the item JSON on stdin and RHIZA_INBOX_KIND/ID/CHANGE/TITLE matching it (e.g. question|d22289a|CSQUESTION|'CSQUESTION: a question waits on you: Which store?'), cwd the repo root, and daemon.state.json recorded hooked ids. A second --once ran none; removing the failed journal then restoring it ran the hook once more. With only git config rhiza.notifyCommand set it ran 3 times; --notify-command fail.sh overrode the config, logged 'notify command for item <id>: exit status 1: boom' 3 times, exit 0, and a rerun logged nothing. --notify-command 'sleep 30' with the real 10 s timeout was killed and logged 'timed out after 10s' per item (30 s total, no sleep left behind); RHIZA_TEST_NOTIFY_TIMEOUT=200ms without RHIZA_TEST_NOW was ignored (still 10 s). --help and --json show notify_command. Not exercised: with Rex attached (would notify the owner). Seen: a hook that backgrounds a child and exits 0 is logged as 'WaitDelay expired before I/O complete' (reported as a minor finding).

## Links

### Out

- verifies [inbox-notify-hook](../capability/rhiza/orchestration/notifications/inbox-notify-hook.md)

## History

- `01M6BD9HT87Z8WCD31EKQQ47R1` The daemon runs an owner-configured command once per new inbox item, with or without Rex
