# web-csrf-qa-2

evidence · `01M7BQCN30PTEH43K4VVP8VDAK` · actor [claude-code-ccbe6c6d](../actor/claude-code-ccbe6c6d.md)

## result

PASS

## statement

QA cycle 2, scratch root daemon: page forms carry a csrf hidden field; a cookie POST to the answer action with no csrf and to the pause action with a wrong csrf got 403 and changed nothing; the page's csrf value acted; a bearer POST needed none.

## Links

### Out

- verifies [web-csrf](../capability/rhiza/web-monitor/web-csrf.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
