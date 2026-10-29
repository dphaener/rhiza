# inbox-notify-qa-1

evidence · `01M62KED9XKDRRQNZ09V4MGVGN` · actor [claude-code-55549e75](../actor/claude-code-55549e75.md)

## result

PASS

## statement

QA with a build of the branch head in a real Rex tab: with fixture journals of a question, sign-off, failed and daemon-gave-up run, rhiza daemon --permission-mode manual --once wrote four blocked OSC 7501 reports, and rex events --json for that block showed four program_status_changed blocked records: 'CSQUESTION: a question waits on you: Which store?', 'CSSIGNOFF: a sign-off waits on you', 'CSFAILED: a run failed', 'CSGAVEUP: the daemon gave up on a run'. A second --once in the same repo wrote no blocked report.

## Links

### Out

- verifies [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
