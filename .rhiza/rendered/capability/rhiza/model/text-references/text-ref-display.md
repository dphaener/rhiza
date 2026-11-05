# text-ref-display

capability · `01M6HP0SAAHDCRXKDTXRXZF7BT` · parent [text-references](../text-references.md)

## acceptance

text-refs-display.txtar: after a node's statement refers to @@other and other is renamed to renamed, rhiza show, why, slice, ls, tree, search, inbox, feed, roadmap, change show and change diff print @@renamed wherever they print that statement, in markdown and in --json fields; rhiza search '@renamed' finds the node; rhiza render writes @@renamed into the rendered page; a reference to a node not in the view shows @@{<full ID>}; after other is retired the reference shows @@{<other's full ID>} and writing the shown statement back with rhiza set stores the same text; a log entry written by hand before the feature with '@word' and '@{ID}' in a statement shows '@@word' and '@@{ID}', and writing that back stores the same literal text; a statement written with '@@word' stores a literal @@word and is shown as @@@word; copying the shown text into rhiza set stores the same text (same IDs, literal @ kept literal)

## statement

Every read operation that prints a node's text fields (render, show, why, slice, ls, tree, search, inbox, feed, roadmap, change show and change diff, the last two in the change set's view), in text and --json alike, show each stored @@{<26-character ID>} as @<the node's current handle> in that view when it names a live node, and as @@{<full ID>} when it names a retired node or does not resolve there; search matches its query against this shown text; and show a stored literal @ that would start a reference as @@@

## Links

### Out

- depends-on [text-ref-syntax](text-ref-syntax.md)
- decided-by [text-references-shape](../../../../decision/text-references-shape.md)

### In

- verifies [text-ref-display-verified](../../../../evidence/text-ref-display-verified.md)
- verifies [text-ref-display-qa-1](../../../../evidence/text-ref-display-qa-1.md)
- verifies [text-ref-display-verified-2](../../../../evidence/text-ref-display-verified-2.md)
- verifies [text-ref-display-qa-2](../../../../evidence/text-ref-display-qa-2.md)
- verifies [pr-278-text-references](../../../../evidence/pr-278-text-references.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
