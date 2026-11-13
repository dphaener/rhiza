# pr-549-fix-rows-on-the-web-roadmap-s

evidence · `01M7AMBV0F43V2Q95C06G8QVTA` · actor [dphaener](../actor/dphaener.md)

## head

6b209449329e2c3cfab855d15a50efbcba03ee8e

## pr

https://github.com/dphaener/rhiza/pull/549

## result

9 checks passed: changes SUCCESS, go-install SUCCESS, model SUCCESS, lint SKIPPED, scripts SKIPPED, api SKIPPED, rest SKIPPED, budgets SKIPPED, check SUCCESS

## statement

Pull request #549 (https://github.com/dphaener/rhiza/pull/549) passed CI at 6b20944

## status

verified

## Links

### Out

- verifies [web-roadmap-rows-resolve](../capability/rhiza/web-monitor/web-roadmap-rows-resolve.md)

## History

- `01M79SMGMW0DCEBFZDHXTDTFEZ` Fix Rows on the web roadmap's backlog sometimes 404 when clicked (owner, 2026-11-13: 'Some of the items in the backlog just 404 when I click on them'). Each row links to /changes/<id> (internal/web/roadmap.go:140), but rows from another project (see the roadmap project-filter bug) or change sets the change page cannot read resolve to no page. Wanted: every roadmap row links to a page that exists in the project it belongs to (under /p/<project>/changes/<id>), and a test clicks every row of a two-project root and gets 200.
