# autonomy-app-verified

evidence · `01M6DM3D1PJN35XGEKDR0JBG3A` · actor [claude-code-daace8aa](../actor/claude-code-daace8aa.md)

## result

PASS

## statement

At 7f3cf020: go test ./cmd/rhiza -run 'TestScripts/(tui-autonomy-app|rex-install)$' passes: --once output cmp-equals `rhiza daemon settings --json` before and after `daemon set rhiza.daemon.maxRuns 3`, and no socket is bound; rex-install greps rhiza_app_autonomy on ["autonomy"] = "y" with requests q, health t, roadmap b, inbox i unmoved. go test ./internal/tui/autonomyapp passes 9 tests on a fake runner: enter toggles rhiza.autonomy.plan owner->auto via daemon.set -- rhiza.autonomy.plan auto and back, enter toggles paused false->true, maxRuns edit keeps the entry open on a refused 0 and saves 5, d+y when not running execs rex window new -s S1 'rhiza daemon' -c /repo --json -- /bin/rhiza daemon --autonomous, n runs nothing, stop sends kill -TERM to the pid on the same host and refuses another host. A scratch repo with a real build showed daemon set -- rhiza.autonomy.plan auto, then --once reporting value auto from git config. go fmt, vet and test in make check pass; self-check passes with the run env unset (with it set, it only flags render files for this run's live actors).

## Links

### Out

- verifies autonomy-app (retired)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
