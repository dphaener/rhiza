# quickstart-plain-commands-qa-1

evidence · `01M73E309JRB5Z7F59JZX7PATE` · actor [claude-code-65f2559d](../actor/claude-code-65f2559d.md)

## result

PASS

## statement

At 8e8c092d I built rhiza with go build -trimpath and walked docs/quickstart.md by hand in a fresh scratch git repo (own HOME, RHIZA_* unset): git init, rhiza init, change new, rhiza add capability reports, add capability export-csv --parent reports, add decision csv-streaming with the quoted quote field, rhiza link export-csv decided-by csv-streaming, propose, accept, land, show, why, tree, render, check (ok), git commit. Every command succeeded. show export-csv lists decided-by csv-streaming, show csv-streaming keeps the quote verbatim, tree nests export-csv under reports, why lists the decision and the landed change set. The page's rhiza set mention works too: set export-csv --field ... --dry-run then real updated the statement. grep finds no {"writes in the page; rhiza apply gets one sentence. go test -trimpath ./cmd/rhiza -run TestScripts/quickstart passes.

## Links

### Out

- verifies [quickstart-plain-commands](../capability/rhiza/documentation/quickstart-plain-commands.md)

## History

- `01M71G4BKFHP1E172RR6H82WFT` README and quickstart in the owner's voice: rewritten with humanize, the quickstart writes the model with plain rhiza add/link commands instead of an apply JSON file, the README's stale 'on the way' lines go, signed off by the owner
