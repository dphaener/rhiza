# import-guide-decision-plan

source · `01M6S9A3Q9XJA62WFXNJA71HGA`

## quotes

```json
[]
```

## raw_sha256

945ba7e9bd9b287cac03201f5f4f3658b3d66c0057239e1a5f6c889edfa8bab8

## session

c7c3877a-453a-4754-9cde-01a1b292e12c

## summary

Autonomous plan run for change set 'The import guide's decision' (no owner present). The planner read rhiza guide import step 4 (internal/guide/text/import.md), the import-autonomous-decisions leaf and its testscript, and found no rhiza check code that reads an imported decision's wording. It raised one question covering three choices: which intent the change set serves, the leaf shape, and what 'rhiza check must accept both forms' means. The arbiter decided (import-owner-credited-wording-shape, assumed, would have escalated): serve milestone-27 and m27-nebari-imported; add one new childless leaf beside import-autonomous-decisions under rhiza/guide, giving step 4 a second wording form '<decided>, made by the owner as recorded in <document>' for decisions a document attributes to the owner, both forms with no quote and a derived-from link to the Source; no new check code, the acceptance being a testscript in which both forms land and pass rhiza check. Nothing is left open.

## title

Plan of the import guide's decision wording, autonomous

## Links

### In

- derived-from [import-owner-credited-wording-shape](../decision/import-owner-credited-wording-shape.md)
- derived-from [import-load-bearing-decisions-owner-credited](../decision/import-load-bearing-decisions-owner-credited.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
