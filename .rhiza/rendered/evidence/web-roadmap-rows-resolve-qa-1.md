# web-roadmap-rows-resolve-qa-1

evidence · `01M7AJF824PCSPM7RGERB435QW` · actor [claude-code-dcf26ddd](../actor/claude-code-dcf26ddd.md)

## result

PASS

## statement

QA at 763883bb6 with a fresh go build -trimpath of ./cmd/rhiza to a scratch path, under an isolated HOME. Made a root (rhiza.toml) with projects alpha and beta, each with a landed milestone intent, an accepted change set serving it, a proposed backlog change set depending on it and a draft; ran the built rhiza daemon (maxRuns 1, autonomy owner, web on 127.0.0.1:17911) and curled the real web UI. Happy: /p/alpha/roadmap and /p/beta/roadmap list only their own change sets (no text from the other project); every row href is under /p/<project>/changes/ and returns 200, and every changes link on those pages (depends-on included) is in the same project and returns 200. Added an umbrella intent with a change set serving it and depending on a landed change set in beta, and registered a third project gamma in rhiza.toml while the daemon ran: the Umbrella group rendered, its row and its depends-on link to the landed change set returned 200, and /p/gamma/roadmap's row returned 200. Every other link on the roadmaps (nav, static, /inbox, /projects, other project's roadmap) returned 200. Failure: /p/alpha/changes/<beta id>, /p/alpha/changes/NOTANID, /p/alpha/changes/ and /p/gamma/* before registration returned a clean 404; old unscoped /changes/<id> and /roadmap redirect 302 to /p/alpha/. /p/alpha/api/roadmap?all=true holds only alpha's change sets. go test -trimpath -count=1 -run TestWebRootRoadmapRowsResolve ./internal/api/ ok; go test -trimpath ./internal/web/... ok.

## Links

### Out

- verifies [web-roadmap-rows-resolve](../capability/rhiza/web-monitor/web-roadmap-rows-resolve.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
