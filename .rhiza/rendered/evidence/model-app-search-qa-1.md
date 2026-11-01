# model-app-search-qa-1

evidence · `01M6BJT5NNQS6GNTY5K699VEYM` · actor [claude-code-85cc7f99](../actor/claude-code-85cc7f99.md)

## result

PASS

## statement

At 7a46ec69, ran the built rhiza tui app model in a real pty (100x30) on a scratch repo and on a clone of this repo's model: the footer shows '/ search'; / opens '/_' with 'esc clear · enter keep · j/k scroll'; typing giz shows only the search hit beta 'Another gizmo'; backspaces and 'plain' re-search to alpha; q types into the filter ('/plainq_', No nodes match.) instead of quitting; on the real model /rex lists rex hits, enter keeps them and j scrolls them, / edits again ('/rexx_'), spaces-only shows no hits; esc restores the tree and its footer, and q then quits with exit 0. rex-install.txtar is untouched by the diff.

## Links

### Out

- verifies [model-app-search](../capability/rhiza/views/text-search/model-app-search.md)

## History

- `01M6BD9GSN00SK8K209A759RT0` rhiza search finds model nodes by their text, and the model app filters with it
