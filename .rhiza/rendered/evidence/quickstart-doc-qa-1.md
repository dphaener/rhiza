# quickstart-doc-qa-1

evidence · `01M6ZATP9S7X2YAXRM3AEX1QB3` · actor [claude-code-d4acd8b3](../actor/claude-code-d4acd8b3.md)

## result

PASS

## statement

QA at b92987bf with a built binary (go build -trimpath ./cmd/rhiza), by hand: in a fresh scratch git repo with its own HOME and every RHIZA_* variable unset, typed each docs/quickstart.md bash block in order (git init, init, change new, apply --dry-run, apply intent.json, propose, accept, land, show/why/slice export-csv, tree, render, check, git add/commit). Every command exited 0; show listed decided-by csv-streaming, why showed the quote and the landed change set, tree showed reports > export-csv, render wrote 5 files, check printed ok, and git status was clean after the commit. go test -trimpath ./cmd/rhiza -run TestScripts/quickstart passes, and fails (unknown command treee, quickstart.txtar:8) after I changed the doc's rhiza tree to rhiza treee, so doc drift is caught.

## Links

### Out

- verifies [quickstart-doc](../capability/rhiza/documentation/quickstart-doc.md)

## History

- `01M6W5EDEX8QAKNVPY5H6CH3Z9` README and quickstart: the README is Rhiza's front door (what and why, install, where next) and docs/quickstart.md walks the model loop without agents, run by a testscript
