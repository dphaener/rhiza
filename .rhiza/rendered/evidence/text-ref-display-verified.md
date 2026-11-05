# text-ref-display-verified

evidence · `01M6KFNFWSBHDM2XE97Y9D80CZ` · actor claude-code-83858632 (retired)

## result

PASS

## statement

At c080ada8 (attempt 2): the attempt-1 finding is fixed: text-refs-write.txtar now expects the shown '@other' (via 49a4311e) and passes. In a clean clone at the commit, rhiza check passes (the self-check failure in the run checkout is only untracked local actor records, mine included) and make fmt-check vet test passes except headless-run-agent-tabs, an unrelated load flake that passes 3/3 alone (go test -count=3 -run 'TestScripts/(headless-run-agent-tabs|text-refs)'). In scratch repos with this build: after @@other was renamed to renamed, show, why, slice, search (text and --json snippets), inbox (via the shared shownAll reader), change show and change diff (text and --json), render and tui app show/slice/why --once all print @@renamed, and search '@renamed' finds the node. ls and tree print no statements; feed and roadmap print only change-set intents. A reference to a retired node, and a hand-written @@{<unknown 26-char ID>}, show as @@{<full ID>}; writing the retired one back with set stores the same @@{ID}. A reference to a node of change set N shows @@newer in N's view, and the as-is view prints the landed text. A hand-edited log entry 'Old @@word and @@{ID}' shows as '@@word and @@@{ID}', and writing it back stores the same literal text. 'Lit @@@word' stores 'Lit @@word' and shows as '@@word'. Read the diff: shown text is applied in the reader (show/read/list/tree/NodesFields), slice, change show/diff, search and render through model.ShowText over textref.Format; source and verbatim quote fields stay raw.

## Links

### Out

- verifies [text-ref-display](../capability/rhiza/model/text-references/text-ref-display.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
