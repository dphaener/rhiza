# import-load-bearing-decisions-owner-credited

decision · `01M6S9G6RVCZZSNXBFMMAVMEV7`

## reasoning

The change set's intent records the owner's pick, credit the owner, for decisions a document attributes to the owner (Nebari's import, nebari#66, had 26). import-load-bearing-decisions worded every imported decision as the agent's and autonomous, which that pick replaces for the owner-attributed ones; the load-bearing rule and the summary of the rest are kept.

## statement

An import brings in as Decision nodes only the decisions that shape a capability or constraint; each one its document records the owner making or confirming is worded '<decided>, made by the owner as recorded in <document>', every other one is worded as the agent's and marked as made autonomously and not confirmed by the owner; none has a quote, each is derived-from its document's Source, and the rest are summarised in the Source of the decisions document.

## wording

(assumed, would have escalated)

## Links

### Out

- supersedes [import-load-bearing-decisions](import-load-bearing-decisions.md)
- derived-from [import-guide-decision-plan](../source/import-guide-decision-plan.md)

### In

- decided-by [import-autonomous-decisions](../capability/rhiza/guide/import-autonomous-decisions.md)
- decided-by [import-owner-credited-decisions](../capability/rhiza/guide/import-owner-credited-decisions.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
