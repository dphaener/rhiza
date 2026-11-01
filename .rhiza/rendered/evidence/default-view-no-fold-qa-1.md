# default-view-no-fold-qa-1

evidence · `01M6C5Z0T73MREAMHA96YW1TRS` · actor [claude-code-0d3d379e](../actor/claude-code-0d3d379e.md)

## result

PASS

## statement

QA at 840de73e with a binary built from the branch and one from the merge-base 608debed. In a scratch repo with a landed, an accepted (root statement changed) and a draft change set (adds draftcap), show root, show draftcap, slice root and tree with no --view gave byte-identical JSON from both binaries with current = landed (as-is), accepted (to-be, 'Root to-be'), draft (what-if:<draft>, draftcap visible), no current (as-is), an unknown ULID (as-is) and garbage in .rhiza/local/change (same error). show root --change by full ID, by prefix, for the landed one, an ambiguous prefix (usage), 'nope' and an unknown ULID (not_found) were identical too. After propose then accept of the current draft, the next read switched to to-be at once (index caught up). On a copy of this repo's .rhiza (5168 log files), show rhiza with a landed current change set and no --view: median 58 ms wall per process on the branch vs 572-632 ms on the base, against 52 ms with no current, 46 ms with --view as-is and 34 ms for rhiza version, with identical output. go test ./internal/api -run 'TestDefaultView|TestReadBudget' -count=1 ok.

## Links

### Out

- verifies [default-view-no-fold](../capability/rhiza/views/fast-reads/default-view-no-fold.md)

## History

- `01M6C3PVXPV5DFQ4G0N6JQGH0Q` A read without --view takes the effective change set's state from the index instead of folding the log
