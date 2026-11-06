# web-overview-daemon-qa-1

evidence · `01M6Q95MANY14Y3VYMQNDXVRGG` · actor [claude-code-885ce195](../actor/claude-code-885ce195.md)

## result

PASS

## statement

At cc2679f4 built rhiza, cloned the repo to a scratch dir, set rhiza.daemon.paused=true, rhiza.daemon.maxRuns=3, rhiza.web.port=18765 and ran the real rhiza daemon. curl http://127.0.0.1:18765/ showed the Daemon panel 'Paused, 0 of 3 slots in use' with State Paused, Run slots 0/3 and three slot dots, and one health pill per part (daemon ok, autonomy warn, machine ok, disk ok, usage ok, main warn, inbox warn, requests ok, roadmap ok, check warn) with the part's state and summary as title. After a poll counted the clone's lock-held runs the heading moved to running slots of 3.

## Links

### Out

- verifies [web-overview-daemon](../capability/rhiza/web-monitor/web-overview-daemon.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
