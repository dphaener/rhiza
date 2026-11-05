# daemon-service-env-verified

evidence · `01M6MH6SZHXYYG9Q0PV1T08S4A` · actor [claude-code-f800dc4a](../actor/claude-code-f800dc4a.md)

## result

PASS

## statement

At 07af60a2 I ran go test -trimpath ./cmd/rhiza -run TestScripts/daemon-install$ -count=1: PASS. The script sets PATH=$WORK/bin:... and TMPDIR='$WORK/tmp dir' (a space), installs, and greps Environment="PATH=...bin..." and Environment="TMPDIR=.../tmp dir" beside ExecStart=... daemon --autonomous; a second install prints 'nothing changed', leaves the unit byte-identical and runs no daemon-reload; a changed PATH ($WORK/extra) prints Installed, rewrites the unit and calls daemon-reload. The diff puts the Environment= lines between WorkingDirectory= and ExecStart=, from ctxGetenv. I also built the binary and ran daemon install in a scratch repo with fake systemctl/loginctl and PATH holding a double quote, a percent, a dollar and a backslash: the unit read Environment="PATH=.../odd\"q%%p$X\\b:..." and Environment="TMPDIR=.../t dir", the whole NAME=value in double quotes with \ and " escaped, % doubled and $ left as is, as systemd reads Environment=. go vet and gofmt are clean on the changed packages.

## Links

### Out

- verifies [daemon-service-env](../capability/rhiza/orchestration/daemon/daemon-service-env.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
