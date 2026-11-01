# whatif-skips-log-scan-qa-1

evidence · `01M6BRE3ZZBFGRHQCM2AXEZ2Q1` · actor [claude-code-6710f3d0](../actor/claude-code-6710f3d0.md)

## result

PASS

## statement

QA at a7cdc22a, hands-on with the built binary against a build of parent c403ceb0. On a clone of the real rhiza repo (1305 change sets, 31 not landed), for every non-landed change set plus a missing ID, ran tree (capability and intent, --retired), ls --retired over all pages, slice --change and show with --view what-if:<cs> under both binaries: byte-identical output (withdrawn ones give the same conflict error, the missing one the same not_found). The whole sweep took 108 s with the parent and 16 s with this build. Running the new binary over an index the parent built (schema upgrade) gave identical output too, as did the change set 01M6BH7Z itself on a clone of the branch head. In a fresh scratch repo, wrote step by step with the new binary (an intention adding a child, setting a landed node and retiring a decision; a second one linking and depending on it; a third with a stale-base conflict once the first landed; a fourth depending on the unlanded third and editing its new node, plus retires) and after each of 9 steps compared every in-flight what-if read against the parent on a fresh copy: identical each time, with the incremental index catch-up in between. Spot-checked: c-node reads cn under what-if:C and cn2 under the dependent D.

## Links

### Out

- verifies [whatif-skips-log-scan](../capability/rhiza/views/fast-reads/whatif-skips-log-scan.md)

## History

- `01M6BH7ZQEKFG0H66J3X50B4TZ` A what-if read folds only the in-flight change sets it needs over the indexed as-is instead of the whole log
