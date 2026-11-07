# import-owner-credited-decisions-verified

evidence · `01M6SN9BRTTPEJ7FPFQ20GSN1X` · actor [claude-code-b6a5ba09](../actor/claude-code-b6a5ba09.md)

## result

PASS

## statement

Verifier claude-code-b6a5ba09 at e4b5d855: read the diff; internal/guide/text/import.md step 4 now gives both wording forms, '<decided>, made by the owner as recorded in <document>' when the document records the owner making or confirming the decision, else the agent-made autonomous form, both with no quote field and a derived-from link to the document's Source (other lines trimmed to stay under the topic token cap). Ran go test -trimpath ./internal/guide/... (ok) and go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|import-.*)$': all 7 pass, including the new import-owner-credited-decisions.txtar, where rhiza guide import prints both forms and when each applies, a repo with one decision of each form, each derived-from its Source and with no quote, lands, passes fold-ok, and rhiza why shows each wording and the Source decisions-md. A fresh checkout's rhiza check fails only on a stale render (index.md plus the builder's committed actor record), which the run's close-out render fixes.

## Links

### Out

- verifies [import-owner-credited-decisions](../capability/rhiza/guide/import-owner-credited-decisions.md)

## History

- `01M6S95EWGHBTY9RSY63NX75AE` The import guide's decision wording is wrong when a document records the owner making or confirming a decision. Step 4 of rhiza guide import forces 'made autonomously by the agent and not confirmed by the owner' for every decision taken from a document; Nebari's import (nebari#66) had 26 decisions its documents attribute to the owner. Owner's pick: credit the owner. The guide must give a second wording form for that case, crediting the owner as the document records it (e.g. '<decided>, made by the owner as recorded in <document>'), still with no quote (they are not the owner's typed words) and a derived-from link to the document's Source; rhiza check must accept both forms.
