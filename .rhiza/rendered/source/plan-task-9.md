# plan-task-9

source · `01M5SWRGWEHYD4B2PN4F734W14`

## quotes

```json
[]
```

## raw_sha256

f2401364ce5381b5662b6e807eb1670270396d658b08fbc0dbefb5097e54bbec

## summary

The shipit work plan for task 9: genesis render, genesis check, genesis ci install, and Genesis checking its own committed .genesis/ in make check. Decided (D1-D14, agent wording accepted by the owner): render writes the as-is view only, an index page plus one page per live node at a handle path, retired nodes listed only in the index; render runs only when asked, and a stale render in check names genesis render as next; check fails on log read problems, fold problems other than concurrent lifecycle entries, as-is problems, link endpoint violations on live as-is links, stale rendered files and append-only violations, and only warns on concurrent entries and on problems that appear only in in-flight views; the append-only base is --base or the merge-base with origin/HEAD, origin/main or main, allowing only added files under .genesis/log and skipping the rule with a warning when there is no base; ci install writes a GitHub workflow and marks rendered/ as generated in .gitattributes, idempotently; Genesis's own .genesis/ is committed and make check runs check with a scratch binary, CI fetching full history; check always returns the full report and exits 1 only on a failure; node page layout; the link endpoint table in one table-tested place; any differing, missing or extra rendered file is stale; tests include golden rendered markdown, a determinism property and a corrupt-log fuzz target; guardrails keep core packages and writes untouched. Anchors and drift stayed out of milestone 1. Nothing left open in the plan.

## title

Shipit plan, milestone 1 task 9: integrity

## Links

### In

- derived-from [check](../capability/rhiza/integrity/check.md)
- derived-from [render](../capability/rhiza/integrity/render.md)
- derived-from [ci-install](../capability/rhiza/integrity/ci-install.md)
- derived-from [self-check](../capability/rhiza/integrity/self-check.md)
- derived-from [t9-d1-render-as-is](../decision/t9-d1-render-as-is.md)
- derived-from [t9-d2-render-explicit](../decision/t9-d2-render-explicit.md)
- derived-from [t9-d3-check-fail-warn](../decision/t9-d3-check-fail-warn.md)
- derived-from [t9-d4-check-base](../decision/t9-d4-check-base.md)
- derived-from [t9-d5-ci-install-files](../decision/t9-d5-ci-install-files.md)
- derived-from [t9-d6-self-check](../decision/t9-d6-self-check.md)
- derived-from [t9-d7-check-report](../decision/t9-d7-check-report.md)
- derived-from [t9-d8-tests](../decision/t9-d8-tests.md)
- derived-from [t9-d9-node-page](../decision/t9-d9-node-page.md)
- derived-from [t9-d10-endpoint-table](../decision/t9-d10-endpoint-table.md)
- derived-from [t9-d11-stale-render](../decision/t9-d11-stale-render.md)
- derived-from [t9-d12-guardrails](../decision/t9-d12-guardrails.md)
- derived-from [t9-d13-defaults](../decision/t9-d13-defaults.md)
- derived-from [t9-d14-acceptance](../decision/t9-d14-acceptance.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
