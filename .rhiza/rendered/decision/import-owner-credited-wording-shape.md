# import-owner-credited-wording-shape

decision · `01M6S997MWKCBVRJ9EQMA5G2RB`

## reasoning

The question names the import-procedure change set's targets, milestone-27 and m27-nebari-imported; milestone-27's statement says what the Nebari import teaches becomes rhiza import, and this change set's intent comes straight from Nebari's import, so those are the fitting landed intents. The existing leaf import-autonomous-decisions specifies only the agent-made form and its testscript; a sibling leaf keeps that leaf's acceptance true and gives the new form its own check, rather than rewriting a built leaf. grep finds no check code reading a decision's wording except the arbiter-answer check in internal/api/run.go, so both forms already pass rhiza check; adding validation the intent did not ask for would be scope creep, and the intent's 'must accept' is satisfied by proving acceptance in a testscript, which also guards against a future check rejecting either form.

## statement

The change set serves milestone-27 and m27-nebari-imported, as the import-procedure change set that built rhiza guide import did: this fixes a gap the Nebari import (nebari#66) exposed, the milestone's own aim of learning from Nebari's import. It is built as one new childless leaf beside import-autonomous-decisions under rhiza/guide (e.g. import-owner-credited-decisions): step 4 of internal/guide/text/import.md gives two wording forms, the existing agent-made one for a document that does not attribute the decision to the owner, and '<decided>, made by the owner as recorded in <document>' for one whose document records the owner making or confirming it; both with no quote field and a derived-from link to the document's Source. import-autonomous-decisions keeps its built agent-made form and testscript, its statement narrowed to decisions the document does not attribute to the owner, and import-load-bearing-decisions is superseded by import-load-bearing-decisions-owner-credited, which states both forms. 'rhiza check must accept both forms' means no new check code: the acceptance is a testscript (extending or beside import-autonomous-decisions.txtar) where rhiza guide import prints both forms, and a repo holding one decision of each form, each derived-from its Source, lands and passes rhiza check (fold-ok), with rhiza why showing each wording and no quote.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c7c3877a-1](../question/question-c7c3877a-1.md)
- derived-from [import-guide-decision-plan](../source/import-guide-decision-plan.md)

### In

- decided-by [import-owner-credited-decisions](../capability/rhiza/guide/import-owner-credited-decisions.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
