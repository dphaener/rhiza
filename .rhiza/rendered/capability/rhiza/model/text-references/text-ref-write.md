# text-ref-write

capability · `01M6HP0GQ76PY04WKDAPVY30Q6` · parent [text-references](../text-references.md)

## acceptance

text-refs-write.txtar: rhiza add capability with statement 'Builds on @@other' stores @@{<other's full ID>} in the log entry; an ID prefix reference stores the full ID; a reference to a node of the write's own change set resolves, one to a node of another unlanded change set or a newly referenced retired node does not; after other is retired, rhiza set of a field that already stores @@{<other's ID>} with input @@{<other's full ID>} keeps it, and adding that reference to a different field or node is refused; quote, quote_where, wording and a source's fields are stored verbatim; '@@' stores a literal @; an input @@@{<full ID>} is refused with code usage (exit 2), naming it and its field, and writes nothing, while @@@<full ID> stores the literal @<full ID>; an input @@{<other's full ID>} stores the same ID, and an input @@{BOGUS} is refused like any unresolvable reference; an unresolvable reference exits 4 with code not_found naming each bad reference and its field, with a next saying to fix the handle or write @@@, and writes nothing; --dry-run reports the same refusal; the same file through rhiza apply behaves the same; rhiza note add with an @@word naming no node is written and shown as @@@word, and one with an @@{<full ID>} naming no node and an @@@{<full ID>} stores each as @<full ID> without braces; describe-golden.txtar is regenerated if a request or response type changes; the rhiza guide's skill text tells agents to write @@handle to refer to a node. Engine writes (decision text-ref-engine-writes): a run session's question (TestAddQuestionLenientRefs) whose text holds an @@word naming no node is written, the run goes on, and it is shown as @@@word, while one holding @<existing handle> stores @@{<full ID>}; text-refs-engine-writes.txtar: question answer behaves the same, its quote kept verbatim; run-ship-change-rebase.txtar: re-raising the same rebase conflict whose statement holds a resolvable reference writes nothing new

## statement

rhiza add and set, and add and set items of rhiza apply, resolve the references in every string-valued field of a node (except quote, quote_where, wording and every field of a source node) at write time and store each as @@{<full ID>}, refusing the whole write when one names no node it may refer to (strict); a reference to a retired node is refused unless the field's current stored value already holds that same ID, in which case it is kept. A strict write of @@@{<full 26-character ID>} is refused with a usage error (code usage, exit 2) naming it and its field, since stored text cannot hold a literal @@{<full ID>} (decision text-ref-stored-escape); @@@<full ID> writes the ID as literal text. Engine writes of free text (a run session's question from rhiza ask, check, leaf and review questions, interview questions, rebase-conflict questions, the Decision written by question answer, notes and ship evidence) resolve the same way but never refuse (lenient): a reference that names no node it may refer to is stored as a literal @, as if written @@@, and shown as @@@word, except that an unresolvable @@{<full ID>} and an @@@{<full ID>} are stored as @<full ID> without braces

## Links

### Out

- depends-on [text-ref-syntax](text-ref-syntax.md)
- decided-by [text-references-shape](../../../../decision/text-references-shape.md)
- decided-by [text-ref-engine-writes](../../../../decision/text-ref-engine-writes.md)
- decided-by [text-ref-stored-escape](../../../../decision/text-ref-stored-escape.md)

### In

- blocks [review-01m6hnnv-r3-f1](../../../../question/review-01m6hnnv-r3-f1.md)
- verifies [text-ref-write-verified](../../../../evidence/text-ref-write-verified.md)
- verifies [text-ref-write-qa-1](../../../../evidence/text-ref-write-qa-1.md)
- verifies [text-ref-write-verified-2](../../../../evidence/text-ref-write-verified-2.md)
- verifies [text-ref-write-qa-2](../../../../evidence/text-ref-write-qa-2.md)
- verifies [pr-278-text-references](../../../../evidence/pr-278-text-references.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
