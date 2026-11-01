# inbox-notify-hook-verified

evidence · `01M6BQK5DF8553ZTG1D8AVS1RY` · actor [claude-code-f1e133da](../actor/claude-code-f1e133da.md)

## result

PASS

## statement

At 060f7afd, go test ./cmd/rhiza -run 'TestScripts/(daemon-notify-hook|daemon-notify|describe-golden)$' -count=1 passed, and make test (gofmt, vet and all packages) passed. daemon-notify-hook.txtar runs without Rex over the four daemon-notify.txtar items: a fake hook appending stdin plus RHIZA_INBOX_* writes 4 lines with question, signoff, failed and gave-up kinds and matching env; a second --once writes none; a new item adds one line; an item removed and returned runs it again; git config rhiza.notifyCommand is used without the flag; fail.sh (exit 1) is logged 4 times with item ID and 'exit status 1: boom', the poll succeeds and it is not rerun; 'sleep 30' under RHIZA_TEST_NOTIFY_TIMEOUT=200ms (honoured only with RHIZA_TEST_NOW, see notifyTimeoutFor) is killed and logged 'timed out after 200ms' 4 times. describe-golden.txtar carries notify_command. The diff keeps Hooked state apart from Notified, so the hook runs with or without Rex. make self-check fails only on stale .rhiza/rendered actor files and index from this run's actor records, not on the leaf's code.

## Links

### Out

- verifies [inbox-notify-hook](../capability/rhiza/orchestration/notifications/inbox-notify-hook.md)

## History

- `01M6BD9HT87Z8WCD31EKQQ47R1` The daemon runs an owner-configured command once per new inbox item, with or without Rex
