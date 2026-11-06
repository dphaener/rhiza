# web-overview-live-qa-2

evidence · `01M6QC68M8G2RWS162JXANE6XX` · actor [claude-code-1f2862fd](../actor/claude-code-1f2862fd.md)

## result

PASS

## statement

At 8166e669 built rhiza (go build -trimpath) to a scratch binary, cloned the branch to a scratch repo with a copy of the real .git/rhiza/runs journals, set rhiza.daemon.paused=true, maxRuns=3, rhiza.web.port=18766 and ran the real rhiza daemon --permission-mode acceptEdits. / declares four live sections: daemon (run daemon), inbox (inbox run change), runs, milestones (change), each with its /fragments/<name> URL, all four 200 alone. With curl -N /events open, writing a new stopped run journal emitted event: run, event: inbox and event: daemon, and the re-fetched inbox fragment then showed the new sign-off row.

## Links

### Out

- verifies [web-overview-live](../capability/rhiza/web-monitor/web-overview-live.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
