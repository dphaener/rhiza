# import-autonomous-decisions

capability · `01M6J3CF3E8B054GW1370AQN2G` · parent [guide](../guide.md)

## acceptance

The import topic states the load-bearing rule and the exact wording form; in a testscript a decision written by that form shows the wording and the Source in rhiza why, and has no quote field.

## statement

The import procedure brings in only decisions that shape a capability or constraint, and an imported decision with no owner words, whose document does not attribute it to the owner, has no quote, a wording naming it as the agent's made autonomously and not confirmed by the owner, and a derived-from link to the Source of the document it came from.

## Links

### Out

- decided-by [import-load-bearing-decisions](../../../decision/import-load-bearing-decisions.md)
- decided-by [nebari-decisions-autonomous](../../../decision/nebari-decisions-autonomous.md)
- decided-by [import-load-bearing-decisions-owner-credited](../../../decision/import-load-bearing-decisions-owner-credited.md)

### In

- verifies [import-autonomous-decisions-verified](../../../evidence/import-autonomous-decisions-verified.md)
- verifies [import-autonomous-decisions-qa-1](../../../evidence/import-autonomous-decisions-qa-1.md)
- verifies [pr-328-import-procedure](../../../evidence/pr-328-import-procedure.md)
- verifies [pr-339-the-import-guide-s-decision](../../../evidence/pr-339-the-import-guide-s-decision.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
