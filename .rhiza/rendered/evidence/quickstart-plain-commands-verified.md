# quickstart-plain-commands-verified

evidence · `01M73D81VZW0320B8S99DCXF4B` · actor [claude-code-f1308b57](../actor/claude-code-f1308b57.md)

## result

PASS

## statement

At 37b4d067, GOTMPDIR= go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/quickstart passes (quickstart-doc ran 19 commands from the page's bash blocks; output shows csv-streaming created and export-csv decided-by csv-streaming). grep finds no {"writes JSON in docs/quickstart.md. The diff replaces the apply heredoc with rhiza add capability (x2), rhiza add decision with quote, and rhiza link export-csv decided-by csv-streaming; rhiza apply gets one short sentence (several writes at once or by an agent). The new Where next link docs/concepts/change-sets.md exists. The commit also stages only its own actor record. Owner approval of the page is left to the run's sign-off.

## Links

### Out

- verifies [quickstart-plain-commands](../capability/rhiza/documentation/quickstart-plain-commands.md)

## History

- `01M71G4BKFHP1E172RR6H82WFT` README and quickstart in the owner's voice: rewritten with humanize, the quickstart writes the model with plain rhiza add/link commands instead of an apply JSON file, the README's stale 'on the way' lines go, signed off by the owner
