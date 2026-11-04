# anchor-fingerprint-shape

decision · `01M6HK5TSPZAT4G3YFQ02XQEKK`

## reasoning

Intent: the umbrella 01M5SWT0K9TV1R serves anyone-builds and model-tells-truth; anchors and fingerprints exist so the model stays true to the code, which is model-tells-truth's statement, so that landed root fits best and anyone-builds is only served indirectly. Fingerprint: the umbrella's anchors capability says path, optional symbol and content fingerprint; sha256 is what the codebase already uses (Source raw_sha256, index and daemon_fold hashing); hashing the declaration rather than the whole file when a symbol is given keeps unrelated edits in the same file from flagging the capability, and Go-only symbols match a Go-only repo without inventing a multi-language parser. Writing: internal/op/fields.go has KindAnchor with no fields, so the three fields are added there; computing the fingerprint at write time means agents never hand-compute hashes, and an explicit re-fingerprint op keeps the fingerprints capability's meaning (a change flags review until someone re-fingerprints) and fits CLAUDE.md's rule that operations live in internal/api behind one registry. Check: the change set's intent says check warns when path or symbol is gone or the fingerprint changed, and check.go already has warn and fail severities; check's paths are relative to the project root, so the working tree is the natural thing to resolve against. implements already links anchor to capability in internal/check/endpoints.go, so findings can name the capability without new link types. Severity: the landed Decision ci-fails-on (owner quote 'yes to all four, move on to A5') says CI fails on broken anchors, so a gone path or symbol fails; the landed fingerprints capability only sanctions a warn for a changed fingerprint.

## statement

The change set serves the landed root intent model-tells-truth. An anchor carries three fields: path (relative to the project root, slash-separated), optional symbol (a top-level Go declaration name, Type.Method for methods; a symbol is only accepted on a .go file) and fingerprint ("sha256:<hex>" of the declaration's source bytes when symbol is set, else of the whole file's bytes, read from the working tree). Writing an anchor (add, or a set of path or symbol) computes the fingerprint when none is given and refuses a path or symbol that does not resolve; one transport-independent operation in internal/api re-fingerprints an existing anchor from the working tree into the current change set, the deliberate step after the code change is reviewed. rhiza check resolves every live anchor in the folded model against the project's working tree, the same checkout CI checks, and reports each finding naming the anchor and the capability it implements: a missing path or a missing symbol is a broken anchor and a fail finding, so rhiza check and CI fail on it as the landed Decision ci-fails-on requires; a changed fingerprint alone is a warn, so editing anchored code never breaks make check or CI until someone reviews it and re-fingerprints.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-298a341f-1](../question/question-298a341f-1.md)
- derived-from [anchors-plan-session](../source/anchors-plan-session.md)

### In

- decided-by [anchors](../capability/rhiza/integrity/anchors.md)
- decided-by [fingerprints](../capability/rhiza/integrity/fingerprints.md)
- decided-by [anchor-fields](../capability/rhiza/integrity/anchors/anchor-fields.md)
- decided-by [check-anchors-resolve](../capability/rhiza/integrity/anchors/check-anchors-resolve.md)
- decided-by [check-fingerprint-changed](../capability/rhiza/integrity/fingerprints/check-fingerprint-changed.md)
- decided-by [anchor-refresh](../capability/rhiza/integrity/fingerprints/anchor-refresh.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
