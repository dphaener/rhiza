# web-overview-inbox-qa-2

evidence · `01M6QC68K1XX0BGBT3DB5Z074E` · actor [claude-code-1f2862fd](../actor/claude-code-1f2862fd.md)

## result

PASS

## statement

At 8166e669 built rhiza (go build -trimpath) to a scratch binary, cloned the branch to a scratch repo with a copy of the real .git/rhiza/runs journals, set rhiza.daemon.paused=true, maxRuns=3, rhiza.web.port=18766 and ran the real rhiza daemon --permission-mode acceptEdits. With only accept items the Inbox showed 'Nothing waits on you'. After adding an open-question stop and a failed open-pr run journal, /fragments/inbox listed '2 waiting on you' in the waiting accent, oldest first: 'Web UI Overview page — open questions wait on your answers' and 'Web UI Roadmap pages — the run failed at its open-pr step', each linking to /runs/<id>; no run ID, question handle, rhiza command or Rex window in any row's text (rhiza inbox's details carried all of them). A later sign-off stop showed 'waiting for your sign-off'.

## Links

### Out

- verifies [web-overview-inbox](../capability/rhiza/web-monitor/web-overview-inbox.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
