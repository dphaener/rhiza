# daemon-service-env-qa-1

evidence · `01M6MHDQQPQ7Q0RW6PK78XMGF7` · actor [claude-code-956b3a04](../actor/claude-code-956b3a04.md)

## result

PASS

## statement

Built rhiza from the branch head and ran `rhiza daemon install` in a scratch repo (scratch HOME, fake systemctl/loginctl on PATH, TMPDIR holding a space). The unit carried Environment="PATH=..." and Environment="TMPDIR=.../tmp dir" above ExecStart=... daemon --autonomous; a second install with the same env printed 'nothing changed' and did not daemon-reload; adding a dir to PATH rewrote the unit and called daemon-reload. A TMPDIR of `t "q" 100% \b $X` was written as `"TMPDIR=t \"q\" 100%% \\b $X"`; loading those Environment= lines into a throwaway real systemd user unit, systemctl --user show read back exactly the original value. daemon-install testscript passes.

## Links

### Out

- verifies [daemon-service-env](../capability/rhiza/orchestration/daemon/daemon-service-env.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
