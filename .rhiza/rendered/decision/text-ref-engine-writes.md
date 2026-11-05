# text-ref-engine-writes

decision · `01M6HPW625P2Y1J7MZSA8YQX4P`

## reasoning

The finding is right that text-references-shape (3) and text-ref-write disagree about engine writes, and both readings it offers break something. If only the CLI requests resolve, engine-written @@handle text is stored literal and shown as @@@handle, so the model no longer tells the truth about which nodes a question names, against model-tells-truth and references-by-id ('stored by ID and shown by handle'). If buildAdd refuses strictly everywhere, an ask or a note becomes a way to fail a run: this plan's own question question-ed7e1e54-1 asks about '@handle names no node', and questions, notes and answers carry prose quoting code, emails and hypothetical handles that an agent or the engine cannot be told to escape. Refusing is for writes whose author is shaping model nodes and can fix a handle (add, set, apply); a question or note exists to record something, and raising or answering it must not fail on wording, which is why rhiza ask and the arbiter flow must always succeed. Storing an unresolvable @@word as a literal @ loses nothing: it is exactly what @@@word would store, it round-trips through display, and rhiza check has nothing dangling to flag. Putting both modes in the one write path keeps one resolver (text-ref-syntax already returns unresolved references with offsets, so lenient mode only escapes them) instead of a second copy per caller. The ship_rebase.go:422 dedupe compares stored text to raw input, which stops matching once a reference is rewritten, so it must compare like with like by resolving the input first.

## statement

Resolution lives in the shared write path (buildAdd and buildSet), so every node write resolves references the same way and nothing written through it is later shown as @@@word. It has two modes. Strict: the generic node writes rhiza add, rhiza set and the add and set items of rhiza apply refuse the whole write on an unresolvable reference, as text-references-shape (3) says. Lenient: every other write of free text, that is a session's question (rhiza ask, run.go), check, leaf and review questions, interview questions, rebase-conflict questions, the Decision written by question.answer and notes, resolves each reference that names a node it may refer to and stores it as @@{<full ID>}, and stores each one that does not resolve as a literal @ (as if written @@@), never refusing; the same verbatim-field exclusions apply in both modes, so question.answer's quote stays exact while its statement is resolved. The rebase-conflict question's dedupe (ship_rebase.go) compares the stored statement with the input after the same lenient resolution, not with the raw input. text-ref-write's acceptance adds: a run session's question whose text holds an @@word naming no node is written, the run goes on, and it is shown as @@@word; one holding @<existing handle> stores @@{<full ID>}; rhiza note and question answer behave the same; and re-raising the same rebase conflict whose statement holds a resolvable reference writes nothing new.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6hnnv-r3-f1](../question/review-01m6hnnv-r3-f1.md)

### In

- decided-by [text-ref-write](../capability/rhiza/model/text-references/text-ref-write.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
