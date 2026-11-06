# daemon-repo-session-qa-1

evidence · `01M6FN74ZZXXA5T17DQ8D1D60K` · actor [claude-code-e6cae1c1](../actor/claude-code-e6cae1c1.md)

## result

PASS

## statement

QA at 78a3b478 with a binary built from the branch, in a scratch repo qarepo with an accepted change set serving a milestone. rhiza daemon --once ran against the real Rex server through a proxy rex that passed session ls/new to real rex and only logged window new. With REX_SESSION=session:other and rhiza-qarepo already present: one 'rex session ls --json', no session new, and the run start tab logged as 'rex window new -s rhiza-qarepo ... run start <cs> --permission-mode acceptEdits'. After rex kill rhiza-qarepo and with REX_SESSION unset: 'rex session new rhiza-qarepo -c <root> --json' ran once, the real server listed rhiza-qarepo afterwards, and the run tab again carried -s rhiza-qarepo. rhiza daemon start (tab mode) also created rhiza-qarepo on the real server and reused it on later starts. When the proxy failed on session ls, the daemon reported the error and opened no tab.

## Links

### Out

- verifies daemon-repo-session (retired)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
