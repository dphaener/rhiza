# web-csrf-qa-1

evidence · `01M7BK2A31J62PRJT02SZ7AYRD` · actor [claude-code-a17d514a](../actor/claude-code-a17d514a.md)

## result

PASS

## statement

QA cycle 1 with the built binary as a scratch root daemon: pages fetched with the token cookie carry a csrf hidden field in every pause, answer and approve form, never the token. Cookie POSTs to /p/proj/actions/pause, /actions/approve and /actions/answer with no csrf, a wrong csrf or the csrf only in the query string got 403 'a valid csrf field is required' and changed nothing; the same POSTs with the page's csrf value acted. A bearer-header POST needed no csrf.

## Links

### Out

- verifies [web-csrf](../capability/rhiza/web-monitor/web-csrf.md)

## History

- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
