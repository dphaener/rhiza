# web-action-pause-qa-1

evidence · `01M7BK2AA7H7225R8JV85GQKVA` · actor [claude-code-a17d514a](../actor/claude-code-a17d514a.md)

## result

PASS

## statement

QA cycle 1 with the built binary as a scratch root daemon: with the cookie and csrf, POST /p/proj/actions/pause scope=project/daemon paused=true answered 'This project is paused.' / 'The daemon is paused.' and wrote rhiza.project.paused=true (project git config) and rhiza.daemon.paused=true (root rhiza.toml); bearer POSTs with paused=false resumed both. The panel's buttons flip to Resume once the daemon polls. scope=world and paused=yes got 'Nothing to pause' / 'Say whether to pause or resume.' Bearer POSTs to /retry, /abandon, /merge, /actions/{retry,abandon,merge}, /p/proj/actions/{retry,merge,abandon}, /p/proj/runs/X/retry, /inbox were 404; /api/daemon.set, /p/proj/api/{daemon.set,run.resume,run.abandon,question.answer} 405; PUT/DELETE/GET on the pause path 405.

## Links

### Out

- verifies [web-action-pause](../capability/rhiza/web-monitor/web-action-pause.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
