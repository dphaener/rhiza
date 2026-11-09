# text-ref-syntax

capability · `01M6HP0CQRJFFH7Q40JCV9H2DG` · parent [text-references](../text-references.md)

## acceptance

Unit tests in internal/textref: a reference starts with @ at the start of the text or after whitespace or an opening ( or [ and is followed by a handle (lower-case words joined by single dashes) or an ID or 4+ character ID prefix; ada@example.com, checkout@v7 and '@key' contain no reference; @@@ yields a literal @ and is never a reference; the resolver returns @@@{<full ID>} as unresolved and flagged escaped and an @@{<full ID>} the lookup does not map as unresolved, writes both as @<full ID> without braces so the formatter shows neither as a reference, and still turns @@@{ID} into the literal @@{ID}; a resolver given a lookup rewrites every reference, @@{<ID>} included, to @@{<full ID>} and returns each unresolved one with its offset, so @@{BOGUS} is unresolved; a formatter treats only @@{<26-character ID>} in stored text as a reference, shows it as @<handle> when the lookup gives a live node's handle and as @@{<full ID>} when the node is retired or the lookup has none, and treats stored '@{ID}', '@{BOGUS}' and '@{<full ID>}' as literal text shown as '@@{ID}' etc.; and shows a stored literal @ that would start a reference (from @@@word) as @@@; a fuzz test checks that formatting resolved text, literal @ included, and resolving it again yields the same stored text

## statement

One pure package, internal/textref, finds the references in a text and turns them both ways: @@handle, @<ID or 4+ character ID prefix> and @@{<ID>} to stored @@{<full ID>}, and stored @@{<full 26-character ID>} back to @<current handle> of a live node (or @@{<full ID>} for a retired or missing one); any other @@{...} in stored text is literal, with @@@ for a literal @ both ways (the formatter re-escapes a literal @ that would start a reference), except that stored text never holds a literal @@{<full 26-character ID>} (decision text-ref-stored-escape): the resolver returns an @@@{<full ID>} escape among the unresolved, flagged as escaped, returns an @@{<full ID>} the lookup does not map as unresolved, and writes both as @<full ID> without braces

## Links

### Out

- decided-by [text-references-shape](../../../../decision/text-references-shape.md)
- decided-by [text-ref-stored-escape](../../../../decision/text-ref-stored-escape.md)

### In

- depends-on [text-ref-write](text-ref-write.md)
- depends-on [text-ref-display](text-ref-display.md)
- depends-on [text-ref-check](text-ref-check.md)
- verifies [text-ref-syntax-verified](../../../../evidence/text-ref-syntax-verified.md)
- verifies [text-ref-syntax-qa-1](../../../../evidence/text-ref-syntax-qa-1.md)
- verifies [text-ref-syntax-verified-2](../../../../evidence/text-ref-syntax-verified-2.md)
- verifies [text-ref-syntax-qa-2](../../../../evidence/text-ref-syntax-qa-2.md)
- verifies [text-ref-syntax-verified-3](../../../../evidence/text-ref-syntax-verified-3.md)
- verifies [pr-278-text-references](../../../../evidence/pr-278-text-references.md)
- verifies [pr-414-pre-publish-audit-and-mit](../../../../evidence/pr-414-pre-publish-audit-and-mit.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
- `01M6W5F0ZPYJ050A3826ARE92T` Pre-publish audit and MIT license: audit git history and .rhiza/ for secrets, personal data and private hostnames, add the MIT LICENSE, and hand the owner the go-ahead to make the repository public and tag v0.1.0
