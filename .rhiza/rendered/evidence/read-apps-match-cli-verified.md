# read-apps-match-cli-verified

evidence · `01M62YPGVYQBJTNTSYVT8JX2FM` · actor [claude-code-74815970](../actor/claude-code-74815970.md)

## result

PASS

## statement

At 608a5c3 ran go test ./cmd/rhiza -run TestScripts/tui-read-apps -count=1 -v: PASS. On a fixture model (init, intent, capability, accepted change set) the script runs rhiza <op> --json and rhiza tui app <op> --once for roadmap, ready, show store, why store and slice store and cmp's the two stdouts; all five match and each --json output is non-empty. --once prints readapps.Model.Result(), the result the app built by running the op through api.Lookup in-process and canon.Marshal. The script asserts no .git/rhiza/daemon.lock or daemon.state.json (daemon.Dir is <git common>/rhiza) and find $WORK -type s prints nothing, with REX_SESSION/REX_SERVER empty. Read the diff of bb5a6c5 and 608a5c3: no daemon or listener is started in the app path. go test ./internal/tui/... and make check pass.

## Links

### Out

- verifies read-apps-match-cli (retired)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
