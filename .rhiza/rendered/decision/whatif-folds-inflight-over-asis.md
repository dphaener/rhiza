# whatif-folds-inflight-over-asis

decision · `01M6BH8H3QZH284QXNYAFZ1K3F`

## statement

A what-if read no longer folds the whole log: it decodes only the in-flight change sets' entries (and the landed ones Revised needs), seeds the model from the index's stored as-is rows, and layers the change set and its closure, giving exactly the full fold's model and problems

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- relates-to [t5-d5-stored-views](t5-d5-stored-views.md) — supersedes it for the what-if part: a what-if no longer runs model.WhatIf over every stored folded change set; it seeds model.WhatIfFrom from the stored as-is rows and decodes only the change set, its in-flight closure and its landed dependencies. As-is and to-be stay materialized as t5-d5 says
- derived-from [m15-orchestrator-picks](../source/m15-orchestrator-picks.md)

### In

- decided-by [whatif-skips-log-scan](../capability/rhiza/views/fast-reads/whatif-skips-log-scan.md)

## History

- `01M6BH7V2AFN6B2PYHQ4NFY3AG` Plan milestone 15: the feed shows the work, what-if reads stop rereading the log, picked by the orchestrator under the owner's standing instruction
