# roadmap-app-start-runs-verified

evidence · `01M675K29J179RQM92Q0BDHBYW` · actor [claude-code-1e668f4b](../actor/claude-code-1e668f4b.md)

## result

PASS

## statement

At 438da7f (attempt 2) I ran go test -count=1 ./internal/tui/roadmapapp -run 'Plan|Ship|Declines|OFocuses' -v and go test -count=1 ./internal/api -run TestRunOpen -v: all pass, and make check exits 0. With a fake exec, p then y makes exactly one rex call 'window new -s session:1 <cs[:9]> -c /repo --json -- /bin/rhiza run start <cs> --workflow plan' with no --autonomous, nothing before y, and the view says 'Opened Rex tab <name>' while staying on the list; s then y ends in 'run start <cs> --autonomous --permission-mode auto'; n after p or s makes no calls. o shows the focused tab or 'No live run'. TestRunOpenFocusesLiveRunTab (new in the fix for finding f1) runs run.open with a fake ctx.Exec: a live run with a recorded block makes 'rex block inspect -s session:s --json block:7' then 'rex focus -s session:s block:7'; an ended run or a change set with no run says not live and calls nothing; an unrecorded or gone block reports tab_gone and does not focus. No test runs real Rex: every rex call goes through an injected exec.

## Links

### Out

- verifies roadmap-app-start-runs (retired)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
