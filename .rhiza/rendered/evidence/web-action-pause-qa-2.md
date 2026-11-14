# web-action-pause-qa-2

evidence · `01M7BQCMXN0KFERP4TBZG8TA25` · actor [claude-code-ccbe6c6d](../actor/claude-code-ccbe6c6d.md)

## result

PASS

## statement

QA cycle 2, scratch root daemon: cookie+csrf POSTs to /p/proj/actions/pause with scope project/daemon and paused true/false answered 'This project is paused.', 'The daemon is paused.', and the resumed forms, leaving rhiza.daemon.paused=false in rhiza.toml and rhiza.project.paused=false in the project's git config. Bearer POSTs to /retry, /abandon, /merge, /actions/retry, /actions/merge, /p/proj/actions/abandon, /p/proj/runs/X/retry were 404; /api/daemon.set, /p/proj/api/run.resume, question.answer and run.approve were 405.

## Links

### Out

- verifies [web-action-pause](../capability/rhiza/web-monitor/web-action-pause.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
