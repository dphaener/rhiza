# pr-339-the-import-guide-s-decision

evidence · `01M6SNVWXRP3KXJ196AYQC5XJP` · actor [dphaener](../actor/dphaener.md)

## head

985ed64372b6ec18b36d1e18afcb57c2fc262921

## pr

https://github.com/dphaener/rhiza/pull/339

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #339 (https://github.com/dphaener/rhiza/pull/339) passed CI at 985ed64

## status

verified

## Links

### Out

- verifies [import-owner-credited-decisions](../capability/rhiza/guide/import-owner-credited-decisions.md)
- verifies [import-autonomous-decisions](../capability/rhiza/guide/import-autonomous-decisions.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
