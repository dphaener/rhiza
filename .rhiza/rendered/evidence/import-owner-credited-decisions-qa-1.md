# import-owner-credited-decisions-qa-1

evidence · `01M6SND2371ZKM1BJFZMG85C1X` · actor [claude-code-b6c050a4](../actor/claude-code-b6c050a4.md)

## result

PASS

## statement

QA at the branch head: built rhiza with go build -trimpath; rhiza guide import step 4 prints both wording forms ('<decided>, made by the owner as recorded in <document>' when the document records the owner making or confirming it, else the agent-made autonomous form), both with no quote field and a derived-from link to the Source. In a fresh scratch git repo, following the guide by hand, added a capability, a decisions-md Source and one decision of each form, each decided-by the capability and derived-from the Source; proposed, accepted, landed and rendered. rhiza check printed ok (exit 0); rhiza why on the capability and on each decision showed each wording and the Source decisions-md; rhiza show --json of each decision had no quote. go test ./internal/guide/... and TestScripts/(guide|import-.*) pass.

## Links

### Out

- verifies [import-owner-credited-decisions](../capability/rhiza/guide/import-owner-credited-decisions.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
