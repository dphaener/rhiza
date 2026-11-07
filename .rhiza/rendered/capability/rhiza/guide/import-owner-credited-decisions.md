# import-owner-credited-decisions

capability · `01M6S9A3X180Y1D4QEHZ8AKDD9` · parent [guide](../guide.md)

## acceptance

In a testscript (beside or extending import-autonomous-decisions.txtar), rhiza guide import prints both wording forms and when each applies; a repo holding one decision of each form, each derived-from its Source and with no quote field, lands and passes rhiza check (fold-ok), and rhiza why shows each decision's wording and its Source.

## statement

The import procedure's step 4 gives a second wording form for a decision its document records the owner making or confirming, '<decided>, made by the owner as recorded in <document>', with no quote field and a derived-from link to the document's Source, beside the agent-made form for every other imported decision.

## Links

### Out

- decided-by [import-owner-credited-wording-shape](../../../decision/import-owner-credited-wording-shape.md)
- decided-by [import-load-bearing-decisions-owner-credited](../../../decision/import-load-bearing-decisions-owner-credited.md)

### In

- verifies [import-owner-credited-decisions-verified](../../../evidence/import-owner-credited-decisions-verified.md)
- verifies [import-owner-credited-decisions-qa-1](../../../evidence/import-owner-credited-decisions-qa-1.md)
- verifies [pr-339-the-import-guide-s-decision](../../../evidence/pr-339-the-import-guide-s-decision.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
